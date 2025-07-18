# Configure the Google Cloud Provider
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }
  required_version = ">= 1.6"
}

# Configure the Google Cloud Provider
provider "google" {
  project = var.project_id
  region  = var.region
}

# Variables
variable "project_id" {
  description = "The GCP project ID"
  type        = string
}

variable "region" {
  description = "The GCP region"
  type        = string
  default     = "us-central1"
}

variable "database_password" {
  description = "Password for the database"
  type        = string
  sensitive   = true
}

variable "jwt_secret" {
  description = "JWT secret for authentication"
  type        = string
  sensitive   = true
}

# Enable required APIs
resource "google_project_service" "apis" {
  for_each = toset([
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "container.googleapis.com",
    "cloudbuild.googleapis.com",
    "run.googleapis.com",
    "sql-component.googleapis.com",
    "sqladmin.googleapis.com",
    "secretmanager.googleapis.com",
    "vpcaccess.googleapis.com",
  ])

  service                    = each.value
  disable_dependent_services = true
  disable_on_destroy         = false
}

# Create a VPC network
resource "google_compute_network" "vpc" {
  name                    = "urlshortener-vpc"
  auto_create_subnetworks = false
  depends_on              = [google_project_service.apis]
}

# Create a subnet
resource "google_compute_subnetwork" "subnet" {
  name          = "urlshortener-subnet"
  ip_cidr_range = "10.0.0.0/24"
  region        = var.region
  network       = google_compute_network.vpc.id
}

# Create VPC Access Connector for Cloud Run
resource "google_vpc_access_connector" "connector" {
  name          = "urlshortener-connector"
  region        = var.region
  ip_cidr_range = "10.8.0.0/28"
  network       = google_compute_network.vpc.id
  depends_on    = [google_project_service.apis]
}

# Create Cloud SQL instance
resource "google_sql_database_instance" "postgres" {
  name             = "urlshortener-db"
  database_version = "POSTGRES_15"
  region           = var.region
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    
    disk_autoresize = true
    disk_size       = 20
    disk_type       = "PD_SSD"

    backup_configuration {
      enabled                        = true
      start_time                     = "02:00"
      point_in_time_recovery_enabled = true
    }

    ip_configuration {
      ipv4_enabled                                  = true
      private_network                               = google_compute_network.vpc.id
      enable_private_path_for_google_cloud_services = true
    }

    database_flags {
      name  = "log_statement"
      value = "all"
    }
  }

  depends_on = [
    google_project_service.apis,
    google_compute_network.vpc
  ]
}

# Create database
resource "google_sql_database" "database" {
  name     = "urlshortener"
  instance = google_sql_database_instance.postgres.name
}

# Create database user
resource "google_sql_user" "user" {
  name     = "urlshortener"
  instance = google_sql_database_instance.postgres.name
  password = var.database_password
}

# Store database connection string in Secret Manager
resource "google_secret_manager_secret" "database_url" {
  secret_id = "DATABASE_URL"
  
  replication {
    automatic = true
  }
  
  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "database_url" {
  secret = google_secret_manager_secret.database_url.id
  secret_data = "postgres://${google_sql_user.user.name}:${var.database_password}@${google_sql_database_instance.postgres.private_ip_address}:5432/${google_sql_database.database.name}?sslmode=require"
}

# Store JWT secret in Secret Manager
resource "google_secret_manager_secret" "jwt_secret" {
  secret_id = "JWT_SECRET"
  
  replication {
    automatic = true
  }
  
  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "jwt_secret" {
  secret = google_secret_manager_secret.jwt_secret.id
  secret_data = var.jwt_secret
}

# Create Cloud Run service for backend
resource "google_cloud_run_service" "backend" {
  name     = "urlshortener-backend"
  location = var.region

  template {
    spec {
      containers {
        image = "gcr.io/${var.project_id}/urlshortener-backend:latest"
        ports {
          container_port = 8080
        }
        
        env {
          name  = "PORT"
          value = "8080"
        }
        
        env {
          name  = "ENVIRONMENT"
          value = "production"
        }
        
        env {
          name = "DATABASE_URL"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.database_url.secret_id
              key  = "latest"
            }
          }
        }
        
        env {
          name = "JWT_SECRET"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.jwt_secret.secret_id
              key  = "latest"
            }
          }
        }

        resources {
          limits = {
            memory = "512Mi"
            cpu    = "1000m"
          }
          requests = {
            memory = "256Mi"
            cpu    = "500m"
          }
        }
      }
      
      container_concurrency = 100
    }
    
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "3"
        "autoscaling.knative.dev/minScale"      = "0"
        "run.googleapis.com/cpu-throttling"     = "true"
        "run.googleapis.com/vpc-access-connector" = google_vpc_access_connector.connector.id
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.apis,
    google_sql_database_instance.postgres,
    google_secret_manager_secret_version.database_url,
    google_secret_manager_secret_version.jwt_secret
  ]
}

# Create Cloud Run service for frontend
resource "google_cloud_run_service" "frontend" {
  name     = "urlshortener-frontend"
  location = var.region

  template {
    spec {
      containers {
        image = "gcr.io/${var.project_id}/urlshortener-frontend:latest"
        ports {
          container_port = 80
        }

        resources {
          limits = {
            memory = "256Mi"
            cpu    = "500m"
          }
          requests = {
            memory = "128Mi"
            cpu    = "250m"
          }
        }
      }
      
      container_concurrency = 100
    }
    
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "2"
        "autoscaling.knative.dev/minScale" = "0"
        "run.googleapis.com/cpu-throttling" = "true"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [google_project_service.apis]
}

# Allow unauthenticated access to Cloud Run services
resource "google_cloud_run_service_iam_member" "backend_noauth" {
  location = google_cloud_run_service.backend.location
  project  = google_cloud_run_service.backend.project
  service  = google_cloud_run_service.backend.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

resource "google_cloud_run_service_iam_member" "frontend_noauth" {
  location = google_cloud_run_service.frontend.location
  project  = google_cloud_run_service.frontend.project
  service  = google_cloud_run_service.frontend.name
  role     = "roles/run.invoker"
  member   = "allUsers"
}

# Outputs
output "backend_url" {
  value = google_cloud_run_service.backend.status[0].url
}

output "frontend_url" {
  value = google_cloud_run_service.frontend.status[0].url
}

output "database_connection_name" {
  value = google_sql_database_instance.postgres.connection_name
}

output "database_private_ip" {
  value = google_sql_database_instance.postgres.private_ip_address
}
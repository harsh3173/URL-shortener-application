#!/bin/bash

# URL Shortener GCP Deployment Script
# This script automates the deployment of the URL shortener to Google Cloud Platform

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
PROJECT_ID=""
REGION="us-central1"
SKIP_TERRAFORM=false
SKIP_BUILD=false

# Function to print colored output
print_message() {
    echo -e "${2}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# Function to check if required tools are installed
check_requirements() {
    print_message "Checking requirements..." $YELLOW
    
    # Check if gcloud is installed
    if ! command -v gcloud &> /dev/null; then
        print_message "gcloud CLI is not installed. Please install it from https://cloud.google.com/sdk/docs/install" $RED
        exit 1
    fi
    
    # Check if terraform is installed
    if ! command -v terraform &> /dev/null; then
        print_message "Terraform is not installed. Please install it from https://terraform.io/downloads.html" $RED
        exit 1
    fi
    
    # Check if docker is installed
    if ! command -v docker &> /dev/null; then
        print_message "Docker is not installed. Please install it from https://docs.docker.com/get-docker/" $RED
        exit 1
    fi
    
    print_message "All requirements met!" $GREEN
}

# Function to authenticate with GCP
authenticate_gcp() {
    print_message "Authenticating with GCP..." $YELLOW
    
    # Check if already authenticated
    if ! gcloud auth list --filter="status:ACTIVE" --format="value(account)" | grep -q .; then
        print_message "Not authenticated with GCP. Running authentication..." $YELLOW
        gcloud auth login
        gcloud auth application-default login
    fi
    
    # Set project if provided
    if [ -n "$PROJECT_ID" ]; then
        gcloud config set project $PROJECT_ID
        print_message "Project set to: $PROJECT_ID" $GREEN
    else
        PROJECT_ID=$(gcloud config get-value project)
        if [ -z "$PROJECT_ID" ]; then
            print_message "No project set. Please set a project with: gcloud config set project PROJECT_ID" $RED
            exit 1
        fi
        print_message "Using current project: $PROJECT_ID" $GREEN
    fi
    
    # Enable required APIs
    print_message "Enabling required APIs..." $YELLOW
    gcloud services enable cloudbuild.googleapis.com
    gcloud services enable run.googleapis.com
    gcloud services enable sql-component.googleapis.com
    gcloud services enable sqladmin.googleapis.com
    gcloud services enable secretmanager.googleapis.com
    gcloud services enable vpcaccess.googleapis.com
}

# Function to build and push Docker images
build_and_push() {
    if [ "$SKIP_BUILD" = true ]; then
        print_message "Skipping build step..." $YELLOW
        return
    fi
    
    print_message "Building and pushing Docker images..." $YELLOW
    
    # Build and push using Cloud Build
    gcloud builds submit --config=deployment/gcp/cloudbuild.yaml .
    
    print_message "Images built and pushed successfully!" $GREEN
}

# Function to deploy infrastructure with Terraform
deploy_infrastructure() {
    if [ "$SKIP_TERRAFORM" = true ]; then
        print_message "Skipping Terraform deployment..." $YELLOW
        return
    fi
    
    print_message "Deploying infrastructure with Terraform..." $YELLOW
    
    cd deployment/gcp/terraform
    
    # Initialize Terraform
    terraform init
    
    # Plan deployment
    terraform plan -var="project_id=$PROJECT_ID" -var="region=$REGION"
    
    # Apply deployment
    print_message "Applying Terraform configuration..." $YELLOW
    terraform apply -var="project_id=$PROJECT_ID" -var="region=$REGION" -auto-approve
    
    # Get outputs
    BACKEND_URL=$(terraform output -raw backend_url)
    FRONTEND_URL=$(terraform output -raw frontend_url)
    
    cd ../../..
    
    print_message "Infrastructure deployed successfully!" $GREEN
    print_message "Backend URL: $BACKEND_URL" $GREEN
    print_message "Frontend URL: $FRONTEND_URL" $GREEN
}

# Function to run post-deployment checks
post_deployment_checks() {
    print_message "Running post-deployment checks..." $YELLOW
    
    # Check if services are responding
    if [ -n "$BACKEND_URL" ]; then
        print_message "Checking backend health..." $YELLOW
        if curl -f "$BACKEND_URL/health" > /dev/null 2>&1; then
            print_message "Backend is healthy!" $GREEN
        else
            print_message "Backend health check failed!" $RED
        fi
    fi
    
    if [ -n "$FRONTEND_URL" ]; then
        print_message "Checking frontend..." $YELLOW
        if curl -f "$FRONTEND_URL" > /dev/null 2>&1; then
            print_message "Frontend is accessible!" $GREEN
        else
            print_message "Frontend check failed!" $RED
        fi
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -p, --project-id PROJECT_ID    GCP Project ID"
    echo "  -r, --region REGION            GCP Region (default: us-central1)"
    echo "  --skip-terraform               Skip Terraform deployment"
    echo "  --skip-build                   Skip Docker build and push"
    echo "  -h, --help                     Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 -p my-project-id"
    echo "  $0 -p my-project-id -r us-east1"
    echo "  $0 -p my-project-id --skip-terraform"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--project-id)
            PROJECT_ID="$2"
            shift 2
            ;;
        -r|--region)
            REGION="$2"
            shift 2
            ;;
        --skip-terraform)
            SKIP_TERRAFORM=true
            shift
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            print_message "Unknown option: $1" $RED
            show_usage
            exit 1
            ;;
    esac
done

# Main deployment flow
main() {
    print_message "Starting URL Shortener deployment to GCP..." $GREEN
    print_message "Project ID: $PROJECT_ID" $YELLOW
    print_message "Region: $REGION" $YELLOW
    
    check_requirements
    authenticate_gcp
    build_and_push
    deploy_infrastructure
    post_deployment_checks
    
    print_message "Deployment completed successfully!" $GREEN
    print_message "Your URL Shortener is now running on Google Cloud Platform!" $GREEN
}

# Run main function
main
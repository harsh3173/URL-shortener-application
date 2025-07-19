# URL Shortener - Production-Ready Full-Stack Application

[![Build Status](https://img.shields.io/badge/build-passing-brightgreen)](https://github.com/your-username/url-shortener)
[![Go Version](https://img.shields.io/badge/go-1.22-blue)](https://golang.org/)
[![React Version](https://img.shields.io/badge/react-18.2-blue)](https://reactjs.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A production-ready URL shortening service built with Go, React, PostgreSQL, and deployed on Google Cloud Platform. Demonstrates full-stack development best practices, scalable architecture, and cloud-native deployment.

## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)
- [API Documentation](#api-documentation)

## Features

### Core Functionality
- URL shortening with custom aliases
- Analytics dashboard with click tracking
- User authentication and management
- Link management (edit, delete, organize)
- Real-time analytics (devices, browsers, geographic data)
- Automatic link expiration
- Rate limiting and abuse prevention
- Responsive, mobile-first UI

### Security & Performance
- OAuth 2.0 authentication with Google
- Input validation and sanitization
- HTTPS enforcement
- Database optimization with connection pooling
- Redis-based caching

## Architecture

```mermaid
graph TB
    subgraph "Frontend (React + TypeScript)"
        A[React App] --> B[Tailwind CSS]
        A --> C[Axios HTTP Client]
        A --> D[React Router]
        A --> E[React Hook Form]
    end
    
    subgraph "Backend (Go + Fiber)"
        F[Fiber Web Server] --> G[OAuth Middleware]
        F --> H[Rate Limiting]
        F --> I[CORS]
        F --> J[API Handlers]
        J --> K[Business Logic]
        K --> L[Database Layer]
    end
    
    subgraph "Database"
        M[PostgreSQL] --> N[User Table]
        M --> O[URL Table]
        M --> P[Analytics Table]
    end
    
    subgraph "Cloud Infrastructure (GCP)"
        Q[Cloud Run - Frontend] --> R[Cloud CDN]
        S[Cloud Run - Backend] --> T[Cloud SQL]
        S --> U[Secret Manager]
        S --> V[VPC Network]
    end
    
    A --> F
    L --> M
    S --> T
    
    subgraph "CI/CD"
        W[GitHub Actions] --> X[Cloud Build]
        X --> Y[Container Registry]
        Y --> Q
        Y --> S
    end
```

## Tech Stack

### Frontend
- **React 18.2** - Modern UI library with hooks
- **TypeScript** - Type-safe JavaScript
- **Tailwind CSS** - Utility-first CSS framework
- **Vite** - Fast build tool and dev server
- **React Router** - Client-side routing
- **Axios** - HTTP client for API calls
- **React Hook Form** - Form validation and management
- **Recharts** - Data visualization for analytics

### Backend
- **Go 1.22** - High-performance backend language
- **Fiber** - Express-inspired web framework
- **GORM** - ORM for database operations
- **PostgreSQL** - Relational database
- **OAuth 2.0** - Google OAuth authentication
- **Sessions** - Secure HTTP-only cookie sessions
- **Testify** - Testing framework

### DevOps & Cloud
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **Google Cloud Platform** - Cloud infrastructure
- **Cloud Run** - Serverless container platform
- **Cloud SQL** - Managed PostgreSQL
- **Cloud Build** - CI/CD pipeline
- **Terraform** - Infrastructure as Code

## Getting Started

### Prerequisites
- Go 1.22+
- Node.js 18+
- Docker and Docker Compose

### Quick Start with Docker

1. Clone the repository
   ```bash
   git clone https://github.com/your-username/url-shortener.git
   cd url-shortener
   ```

2. Start the development environment
   ```bash
   make dev
   ```

3. Access the application
   - Frontend: http://localhost:3001
   - Backend API: http://localhost:8080
   - Database: postgresql://postgres:password@localhost:5432/urlshortener

### Manual Setup

#### Backend Setup
1. Navigate to backend directory: `cd backend`
2. Install dependencies: `go mod download`
3. Set up environment variables: `cp .env.example .env`
4. Run the server: `go run cmd/main.go`

#### Frontend Setup
1. Navigate to frontend directory: `cd frontend`
2. Install dependencies: `npm install`
3. Set up environment variables: `cp .env.example .env`
4. Start the development server: `npm run dev`

## Development

### Available Commands

```bash
make dev     # Start development environment
make build   # Build all services
make test    # Run tests
make logs    # View logs
make down    # Stop all services
make clean   # Clean up Docker resources
```

### Project Structure

```
url-shortener/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   │   ├── config/        # Configuration management
│   │   ├── database/      # Database connection and migrations
│   │   ├── handlers/      # HTTP handlers
│   │   ├── middleware/    # Custom middleware
│   │   ├── models/        # Data models
│   │   ├── services/      # Business logic
│   │   └── utils/         # Utility functions
│   ├── tests/             # Test files
│   └── Dockerfile         # Docker configuration
├── frontend/               # React frontend application
│   ├── src/
│   │   ├── components/    # Reusable UI components
│   │   ├── pages/         # Page components
│   │   ├── services/      # API services
│   │   ├── types/         # TypeScript type definitions
│   │   └── utils/         # Utility functions
│   ├── public/            # Static assets
│   └── Dockerfile         # Docker configuration
├── deployment/            # Deployment configurations
│   ├── docker/           # Docker-related files
│   └── gcp/              # Google Cloud Platform configs
└── docker-compose.yml     # Multi-container setup
```

## Testing

### Backend Tests
```bash
cd backend && go test ./tests/... -v        # Run all tests
cd backend && go test ./tests/... -v -cover # Run with coverage
```

### Frontend Tests
```bash
cd frontend && npm test                     # Run all tests
cd frontend && npm run test:coverage        # Run with coverage
```

### Integration Tests
```bash
make test           # Run full test suite
make test-backend   # Run backend tests only
make test-frontend  # Run frontend tests only
```

## Deployment

### Local Production Build
```bash
make up  # Build and run production containers
```

### Google Cloud Platform Deployment

#### Prerequisites
- GCP account with billing enabled
- `gcloud` CLI installed and authenticated
- `terraform` installed

#### Automated Deployment
```bash
chmod +x deployment/gcp/deploy.sh
./deployment/gcp/deploy.sh -p your-project-id -r us-central1
```

#### Manual Deployment
1. Set up GCP project
   ```bash
   gcloud config set project your-project-id
   gcloud services enable cloudbuild.googleapis.com run.googleapis.com
   ```

2. Deploy infrastructure
   ```bash
   cd deployment/gcp/terraform
   terraform init
   terraform plan -var="project_id=your-project-id"
   terraform apply -var="project_id=your-project-id"
   ```

3. Build and deploy containers
   ```bash
   gcloud builds submit --config=deployment/gcp/cloudbuild.yaml .
   ```

## API Documentation

### Authentication Endpoints

#### GET /api/v1/auth/login
Initiate OAuth 2.0 login flow with Google.

Returns:
```json
{
  "success": true,
  "data": {
    "auth_url": "https://accounts.google.com/oauth2/auth?..."
  },
  "message": "OAuth login URL generated"
}
```

### URL Management Endpoints

#### POST /api/v1/urls
Create a new shortened URL.

```json
{
  "original_url": "https://example.com/very-long-url",
  "custom_alias": "my-link",
  "title": "My Custom Link",
  "description": "Description of the link",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

#### GET /api/v1/urls
Get user's URLs with pagination.

Query Parameters:
- `limit` (default: 10) - Number of URLs to return
- `offset` (default: 0) - Number of URLs to skip

#### GET /:shortCode
Redirect to original URL and record analytics.

### Analytics Endpoints

#### GET /api/v1/urls/:id/analytics
Get detailed analytics for a specific URL.

Response includes:
- Click statistics
- Device and browser breakdown
- Geographic data
- Time-based analytics

## Technical Highlights

### Backend Development
- Go programming with modern patterns and best practices
- RESTful API design and implementation
- Database design and optimization with PostgreSQL
- Authentication and security (OAuth 2.0, sessions)
- Middleware development and request handling
- Comprehensive testing with unit and integration tests

### Frontend Development
- React 18 with modern hooks and functional components
- TypeScript for type-safe development
- Responsive design with Tailwind CSS
- State management and form handling
- API integration and error handling
- Component-based architecture

### DevOps & Cloud
- Docker containerization and multi-service orchestration
- Google Cloud Platform deployment (Cloud Run, Cloud SQL)
- Infrastructure as Code with Terraform
- CI/CD pipeline setup with Cloud Build
- Security best practices (Secret Manager, VPC)
- Performance optimization and scalability

### Key Achievements
- Full-stack development with complete application from database to UI
- Cloud-native architecture deployed on GCP with auto-scaling capabilities
- Security-first approach with OAuth 2.0 authentication, input validation, and HTTPS
- Test-driven development with comprehensive test suite
- Production-ready with monitoring, logging, and error handling
- Scalable design with horizontal scaling and caching
- Modern tech stack with latest versions of Go, React, and cloud services

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
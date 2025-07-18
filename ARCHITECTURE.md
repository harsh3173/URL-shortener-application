# URL Shortener Architecture Documentation

## System Overview

The URL Shortener is a modern, cloud-native application built with a microservices architecture. It consists of three main components:

1. **React Frontend** - User interface and client-side logic
2. **Go Backend** - API server and business logic
3. **PostgreSQL Database** - Data persistence layer

## Architecture Principles

### 1. Separation of Concerns
- **Frontend**: User interface and user experience
- **Backend**: Business logic and data management
- **Database**: Data persistence and integrity

### 2. Scalability
- Stateless backend services
- Horizontal scaling capability
- Database connection pooling
- Caching layer for performance

### 3. Security
- JWT-based authentication
- Input validation and sanitization
- HTTPS enforcement
- Rate limiting and abuse prevention

### 4. Maintainability
- Clean code architecture
- Comprehensive testing
- Documentation and code comments
- Modular design patterns

## Component Details

### Frontend Architecture

```
src/
├── components/          # Reusable UI components
│   ├── forms/          # Form-specific components
│   ├── ui/             # Basic UI elements
│   └── layout/         # Layout components
├── pages/              # Page-level components
├── services/           # API communication
├── contexts/           # React contexts for state
├── hooks/              # Custom React hooks
├── utils/              # Utility functions
└── types/              # TypeScript type definitions
```

**Key Technologies:**
- **React 18**: Modern UI library with concurrent features
- **TypeScript**: Type safety and better developer experience
- **Tailwind CSS**: Utility-first CSS framework
- **Vite**: Fast build tool and development server
- **React Router**: Client-side routing
- **Axios**: HTTP client for API calls

### Backend Architecture

```
internal/
├── config/             # Configuration management
├── database/           # Database connection and migrations
├── handlers/           # HTTP request handlers
├── middleware/         # Custom middleware
├── models/             # Data models and structures
├── services/           # Business logic layer
└── utils/              # Utility functions
```

**Key Technologies:**
- **Go 1.22**: High-performance, statically typed language
- **Fiber**: Express-inspired web framework
- **GORM**: Object-Relational Mapping for Go
- **JWT**: JSON Web Token authentication
- **bcrypt**: Password hashing
- **PostgreSQL**: Relational database

### Database Schema

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- URLs table
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    custom_alias VARCHAR(50) UNIQUE,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(200),
    description TEXT,
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Analytics table
CREATE TABLE analytics (
    id SERIAL PRIMARY KEY,
    url_id INTEGER REFERENCES urls(id),
    ip_address INET,
    user_agent TEXT,
    referrer TEXT,
    country VARCHAR(100),
    city VARCHAR(100),
    device VARCHAR(100),
    os VARCHAR(100),
    browser VARCHAR(100),
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Data Flow

### URL Creation Flow

1. **User Input**: User submits URL through frontend form
2. **Validation**: Frontend validates input before sending to backend
3. **Authentication**: Backend verifies JWT token (if user is logged in)
4. **Processing**: Backend generates short code and validates URL
5. **Storage**: URL data is stored in PostgreSQL database
6. **Response**: Short URL is returned to frontend
7. **Display**: Frontend displays the shortened URL to user

### URL Redirection Flow

1. **Request**: User clicks on shortened URL
2. **Lookup**: Backend queries database for original URL
3. **Analytics**: Click data is recorded asynchronously
4. **Validation**: Check if URL is active and not expired
5. **Redirect**: HTTP redirect response to original URL
6. **Tracking**: Analytics data is processed and stored

### Analytics Flow

1. **Data Collection**: Each click generates analytics data
2. **Processing**: User agent parsing, geolocation, device detection
3. **Storage**: Analytics data stored in dedicated table
4. **Aggregation**: Real-time aggregation for dashboard display
5. **Visualization**: Charts and graphs rendered in frontend

## Security Architecture

### Authentication Flow

1. **Registration/Login**: User credentials are validated
2. **Password Hashing**: bcrypt hashing with salt
3. **JWT Generation**: Signed JWT token is created
4. **Token Storage**: HttpOnly cookies for security
5. **Token Validation**: Middleware validates tokens on protected routes

### Security Measures

- **Input Validation**: All inputs are validated and sanitized
- **SQL Injection Prevention**: GORM provides automatic escaping
- **XSS Protection**: React's built-in XSS protection
- **CSRF Protection**: SameSite cookie attributes
- **Rate Limiting**: Request throttling to prevent abuse
- **HTTPS Enforcement**: All production traffic encrypted

## Performance Optimizations

### Frontend Optimizations

- **Code Splitting**: Dynamic imports for route-based splitting
- **Lazy Loading**: Components loaded on demand
- **Memoization**: React.memo and useMemo for expensive operations
- **Image Optimization**: Responsive images and lazy loading
- **Bundle Analysis**: Webpack bundle analyzer for optimization

### Backend Optimizations

- **Database Indexing**: Strategic indexes on frequently queried columns
- **Connection Pooling**: Efficient database connection management
- **Query Optimization**: Efficient SQL queries and joins
- **Caching**: Redis caching for frequently accessed data
- **Compression**: Response compression for bandwidth savings

### Database Optimizations

- **Indexing Strategy**:
  ```sql
  CREATE INDEX idx_urls_short_code ON urls(short_code);
  CREATE INDEX idx_urls_user_id ON urls(user_id);
  CREATE INDEX idx_analytics_url_id ON analytics(url_id);
  CREATE INDEX idx_analytics_clicked_at ON analytics(clicked_at);
  ```

- **Query Optimization**: Efficient joins and aggregations
- **Connection Pooling**: Optimized connection pool settings
- **Backup Strategy**: Automated backups with point-in-time recovery

## Deployment Architecture

### Local Development

```
Docker Compose Environment:
├── Frontend (React Dev Server)
├── Backend (Go Application)
├── PostgreSQL Database
├── Redis Cache (Optional)
└── Adminer (Database Management)
```

### Production Deployment (GCP)

```
Google Cloud Platform:
├── Cloud Run (Frontend)
├── Cloud Run (Backend)
├── Cloud SQL (PostgreSQL)
├── Cloud Build (CI/CD)
├── Container Registry
├── Secret Manager
├── VPC Network
└── Load Balancer
```

## Monitoring and Observability

### Logging Strategy

- **Structured Logging**: JSON-formatted logs for parsing
- **Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Request Logging**: HTTP request/response logging
- **Error Tracking**: Detailed error information and stack traces

### Metrics Collection

- **Application Metrics**: Response times, error rates, throughput
- **Business Metrics**: URL creation rates, click-through rates
- **Infrastructure Metrics**: CPU, memory, disk usage
- **Database Metrics**: Query performance, connection pool stats

### Health Checks

- **Application Health**: `/health` endpoint for service status
- **Database Health**: Connection and query validation
- **Dependency Health**: External service availability
- **Load Balancer Health**: Integration with GCP health checks

## Scalability Considerations

### Horizontal Scaling

- **Stateless Services**: No server-side session storage
- **Load Balancing**: Multiple backend instances
- **Database Scaling**: Read replicas and connection pooling
- **Cache Distribution**: Redis cluster for high availability

### Vertical Scaling

- **Resource Optimization**: CPU and memory tuning
- **Database Optimization**: Query performance tuning
- **Connection Limits**: Optimal connection pool sizing
- **Cache Sizing**: Memory allocation for caching

### Future Enhancements

- **CDN Integration**: Global content delivery
- **Multi-region Deployment**: Reduced latency worldwide
- **Event-driven Architecture**: Async processing with pub/sub
- **Microservices Split**: Separate analytics and URL services

## Testing Strategy

### Unit Testing

- **Backend**: Go test framework with mocks
- **Frontend**: Jest and React Testing Library
- **Coverage**: 80%+ code coverage requirement
- **Mocking**: External dependencies mocked

### Integration Testing

- **API Testing**: End-to-end API workflows
- **Database Testing**: CRUD operations validation
- **Authentication Testing**: JWT token validation
- **Error Handling**: Error scenarios and recovery

### End-to-End Testing

- **User Workflows**: Complete user journey testing
- **Cross-browser Testing**: Multiple browser compatibility
- **Mobile Testing**: Responsive design validation
- **Performance Testing**: Load testing and stress testing

This architecture provides a solid foundation for a production-ready URL shortener service with room for future enhancements and scaling.
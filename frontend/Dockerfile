# Build stage
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Add build dependencies
RUN apk add --no-cache git

# Copy package files
COPY package*.json ./

# Install ALL dependencies (dev + prod for build)
RUN npm ci --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM nginx:alpine

# Install security updates
RUN apk upgrade --no-cache

# Create non-root user for security
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001

# Copy built files from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Fix permissions and create directories
RUN chown -R 1001:1001 /usr/share/nginx/html && \
    chown -R 1001:1001 /var/cache/nginx && \
    chown -R 1001:1001 /var/log/nginx && \
    chown -R 1001:1001 /etc/nginx/conf.d && \
    chown -R 1001:1001 /tmp && \
    mkdir -p /var/run/nginx && \
    chown -R 1001:1001 /var/run/nginx && \
    mkdir -p /tmp/nginx && \
    chown -R 1001:1001 /tmp/nginx

# Change default nginx port to 8080 (non-privileged)
RUN sed -i.bak 's/listen\s*80;/listen 8080;/' /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 8080

# Switch to non-root user
USER 1001

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
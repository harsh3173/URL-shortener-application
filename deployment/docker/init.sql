-- Initialize the database with any required setup
-- This file is run when the PostgreSQL container starts for the first time

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create indexes for better performance (GORM will create tables automatically)
-- These will be applied after tables are created by the application

-- Note: The actual table creation is handled by GORM AutoMigrate in the Go application
-- This file is for any additional database setup that might be needed
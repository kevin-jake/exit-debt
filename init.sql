-- PostgreSQL initialization script for debt tracker
-- This file runs automatically when the container starts for the first time

-- Create extensions if needed
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- You can add any initial database setup here
-- For example, creating additional schemas, users, or initial data

-- Example: Create a schema for the application
-- CREATE SCHEMA IF NOT EXISTS pay_your_dues;

-- Example: Create initial tables (uncomment and modify as needed)
/*
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS debts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*/ 
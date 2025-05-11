-- PostgreSQL initialization script

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE wander_sphere'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'wander_sphere');

-- Connect to the wander_sphere database
\c wander_sphere;

-- Set up extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; 
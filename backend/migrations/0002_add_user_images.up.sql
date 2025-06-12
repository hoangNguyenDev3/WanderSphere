-- Add profile and cover picture columns to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS profile_picture VARCHAR(1000),
ADD COLUMN IF NOT EXISTS cover_picture VARCHAR(1000); 
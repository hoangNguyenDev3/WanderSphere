-- Remove profile and cover picture columns from users table
ALTER TABLE users 
DROP COLUMN IF EXISTS profile_picture,
DROP COLUMN IF EXISTS cover_picture; 
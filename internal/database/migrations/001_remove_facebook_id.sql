-- Migration to remove facebook_id column from contacts table
-- This migration should be run when upgrading from a version that had Facebook ID support

-- Remove facebook_id column from contacts table
ALTER TABLE contacts DROP COLUMN IF EXISTS facebook_id;

-- Remove recipient_facebook_id column from notifications table
ALTER TABLE notifications DROP COLUMN IF EXISTS recipient_facebook_id;

-- Add unique constraint to phone column in contacts table
-- Note: This will fail if there are duplicate phone numbers in the database
-- You may need to clean up duplicates first
ALTER TABLE contacts ADD CONSTRAINT contacts_phone_unique UNIQUE (phone);

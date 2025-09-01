-- Migration to remove facebook_id column from contacts table
-- This migration should be run when upgrading from a version that had Facebook ID support

-- Remove facebook_id column from contacts table
ALTER TABLE contacts DROP COLUMN IF EXISTS facebook_id;

-- Remove recipient_facebook_id column from notifications table
ALTER TABLE notifications DROP COLUMN IF EXISTS recipient_facebook_id;

-- Migration to update debt_items status constraint to include 'rejected' status
-- This migration fixes the constraint violation when trying to set status to 'rejected'

-- Drop the existing constraint
ALTER TABLE debt_items DROP CONSTRAINT IF EXISTS chk_debt_items_status;

-- Add the updated constraint with 'rejected' status included
ALTER TABLE debt_items ADD CONSTRAINT chk_debt_items_status 
    CHECK (status IN ('completed', 'pending', 'failed', 'refunded', 'rejected'));




-- Migration to add user_id column to debt_items table
-- This migration adds a user_id column to track which user created each payment

-- Step 1: Add user_id column as nullable first (to handle existing data)
ALTER TABLE debt_items ADD COLUMN IF NOT EXISTS user_id UUID;

-- Step 2: Update existing rows to set user_id from the debt_list
-- All debt items belong to a debt list, and the debt list has a user_id
-- So we can populate user_id from the parent debt_list
UPDATE debt_items 
SET user_id = debt_lists.user_id
FROM debt_lists
WHERE debt_items.debt_list_id = debt_lists.id
  AND debt_items.user_id IS NULL;

-- Step 3: Make the column NOT NULL (now that all rows have values)
ALTER TABLE debt_items ALTER COLUMN user_id SET NOT NULL;

-- Step 4: Add an index for performance on user_id lookups
CREATE INDEX IF NOT EXISTS idx_debt_items_user_id ON debt_items(user_id);

-- Step 5: Optionally add a foreign key constraint to ensure referential integrity
-- ALTER TABLE debt_items ADD CONSTRAINT fk_debt_items_user 
--     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;


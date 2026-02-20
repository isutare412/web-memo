-- DML migration: Backfill publish_state from is_published
-- Run AFTER Ent auto-migrates DDL (which adds the publish_state column)

-- Backfill publish_state based on existing is_published values
UPDATE memos SET publish_state = 'published' WHERE is_published = true;
UPDATE memos SET publish_state = 'private' WHERE is_published = false;

-- Backfill approved for existing subscriptions (all existing subscriptions are approved)
UPDATE subscriptions SET approved = true WHERE approved IS NULL;

-- Drop the old is_published column after backfilling publish_state
ALTER TABLE memos DROP COLUMN is_published;

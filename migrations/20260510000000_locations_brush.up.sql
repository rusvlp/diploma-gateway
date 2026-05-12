ALTER TABLE locations
    ADD COLUMN IF NOT EXISTS extra_trees       TEXT NOT NULL DEFAULT '[]',
    ADD COLUMN IF NOT EXISTS removed_positions TEXT NOT NULL DEFAULT '[]';

ALTER TABLE locations
    DROP COLUMN IF EXISTS extra_trees,
    DROP COLUMN IF EXISTS removed_positions;

-- Drop triggers first
DROP TRIGGER IF EXISTS entity_path_update ON z_entity;
DROP TRIGGER IF EXISTS entity_path_insert ON z_entity;

-- Drop functions
DROP FUNCTION IF EXISTS update_entity_path();

-- Drop indexes
DROP INDEX IF EXISTS idx_entity_type;
DROP INDEX IF EXISTS idx_entity_parent_id;
DROP INDEX IF EXISTS idx_entity_path;

-- Drop tables
DROP TABLE IF EXISTS z_entity;

-- Drop extensions
DROP EXTENSION IF EXISTS ltree;

-- Make sure ltree extension is created before the table that uses it
CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE IF NOT EXISTS z_entity (
  entity_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
  name varchar(255) NOT NULL,
  details JSONB NOT NULL DEFAULT '{}'::jsonb,
  category_id UUID NOT NULL,
  parent_id UUID,
  path LTREE,
  depth INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(parent_id) REFERENCES z_entity(entity_id) ON DELETE SET NULL,
  FOREIGN KEY(category_id) REFERENCES z_category(category_id) ON DELETE CASCADE
);

-- Fixed table name in indexes (from 'entity' to 'z_entity')
CREATE INDEX idx_entity_path ON z_entity USING GIST(path);
CREATE INDEX idx_entity_parent_id ON z_entity(parent_id);

CREATE OR REPLACE FUNCTION update_entity_path() RETURNS TRIGGER AS $$
DECLARE
    parent_path ltree;
    new_depth integer;
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.path = text2ltree(NEW.entity_id::text);
        NEW.depth = 1;
        RETURN NEW;
    END IF;
    
    SELECT path, depth INTO parent_path, new_depth
    FROM z_entity
    WHERE entity_id = NEW.parent_id;

    IF parent_path IS NULL THEN
        RAISE EXCEPTION 'parent entity not found';
    END IF;

    NEW.path = parent_path || text2ltree(NEW.entity_id::text);
    NEW.depth = new_depth + 1;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER entity_path_insert 
BEFORE INSERT ON z_entity
FOR EACH ROW 
EXECUTE FUNCTION update_entity_path();

CREATE TRIGGER entity_path_update 
BEFORE UPDATE ON z_entity 
FOR EACH ROW 
WHEN (OLD.parent_id IS DISTINCT FROM NEW.parent_id) 
EXECUTE FUNCTION update_entity_path();


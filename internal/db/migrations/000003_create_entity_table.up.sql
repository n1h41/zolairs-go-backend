-- Make sure ltree extension is created before the table that uses it
CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE IF NOT EXISTS z_entity (
    entity_id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid,
    name varchar(255) NOT NULL,
    details jsonb NOT NULL DEFAULT '{}'::jsonb,
    category_id uuid NOT NULL,
    parent_id uuid,
    path ltree,
    depth int NOT NULL DEFAULT 1,
    created_at timestamp with time zone DEFAULT current_timestamp,
    updated_at timestamp with time zone DEFAULT current_timestamp,
    FOREIGN KEY (user_id) REFERENCES z_users (user_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES z_entity (entity_id) ON DELETE SET NULL,
    FOREIGN KEY (category_id) REFERENCES z_category (
        category_id
    ) ON DELETE CASCADE
);

-- Fixed table name in indexes (from 'entity' to 'z_entity')
CREATE INDEX idx_entity_path ON z_entity USING gist (path);

CREATE INDEX idx_entity_parent_id ON z_entity (parent_id);

CREATE OR REPLACE FUNCTION update_entity_path()
RETURNS trigger
AS $$
DECLARE
    parent_path ltree;
    new_depth integer;
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.path = text2ltree (NEW.entity_id::text);
        NEW.depth = 1;
        RETURN NEW;
    END IF;
    SELECT
        path,
        depth INTO parent_path,
        new_depth
    FROM
        z_entity
    WHERE
        entity_id = NEW.parent_id;
    IF parent_path IS NULL THEN
        RAISE EXCEPTION 'parent entity not found';
    END IF;
    NEW.path = parent_path || text2ltree (NEW.entity_id::text);
    NEW.depth = new_depth + 1;
    RETURN NEW;
END;
$$
LANGUAGE plpgsql;

CREATE TRIGGER entity_path_insert
BEFORE INSERT ON z_entity
FOR EACH ROW
EXECUTE FUNCTION update_entity_path();

CREATE TRIGGER entity_path_update
BEFORE UPDATE ON z_entity
FOR EACH ROW
WHEN (old.parent_id IS DISTINCT FROM new.parent_id)
EXECUTE FUNCTION update_entity_path();

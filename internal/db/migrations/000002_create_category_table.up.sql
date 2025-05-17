DO $$
BEGIN
  IF NOT EXISTS(SELECT 1 FROM pg_type WHERE typname = 'category_type') THEN
    CREATE TYPE category_type AS ENUM ('user', 'office', 'location');
  END IF;
END
$$;

CREATE TABLE IF NOT EXISTS z_category (
  category_id UUID PRIMARY KEY NOT NULL DEFAULT GEN_RANDOM_UUID(),
  name VARCHAR(255) NOT NULL,
  type category_type NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Drop the table first (since it depends on the type)
DROP TABLE IF EXISTS z_category;

-- Drop the enum type
DO $$
BEGIN
  IF EXISTS(SELECT 1 FROM pg_type WHERE typname = 'category_type') THEN
    DROP TYPE category_type;
  END IF;
END
$$;


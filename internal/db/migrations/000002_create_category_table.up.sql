DO $$
BEGIN
    IF NOT EXISTS (
        SELECT
            1
        FROM
            pg_type
        WHERE
            typname = 'category_type') THEN
    CREATE TYPE category_type AS ENUM (
        'user',
        'office',
        'location'
);
END IF;
END
$$;

CREATE TABLE IF NOT EXISTS z_category (
    category_id uuid PRIMARY KEY NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name varchar(255) NOT NULL,
    type category_type NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT
            1
        FROM
            pg_type
        WHERE
            typname = 'user_role') THEN
    CREATE TYPE user_role AS ENUM (
        'admin',
        'user'
);
END IF;
END
$$;

CREATE TABLE IF NOT EXISTS z_users (
    user_id uuid PRIMARY KEY NOT NULL DEFAULT GEN_RANDOM_UUID(),
    email varchar(255) UNIQUE NOT NULL,
    first_name varchar(100),
    last_name varchar(100),
    phone varchar(50),
    address jsonb,
    parent_id uuid,
    cognito_id varchar(255) NOT NULL UNIQUE,
    role user_role NOT NULL DEFAULT 'user',
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES z_users (user_id)
);

CREATE INDEX idx_users_email ON z_users (email);

CREATE INDEX idx_users_parent_id ON z_users (parent_id);

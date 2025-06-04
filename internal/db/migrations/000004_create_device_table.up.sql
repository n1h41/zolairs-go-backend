CREATE TABLE z_device (
    mac_address varchar(17) PRIMARY KEY NOT NULL UNIQUE,
    user_id uuid NOT NULL,
    device_name varchar(255) NOT NULL,
    category varchar(255),
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES z_users (user_id) ON DELETE CASCADE
);

CREATE INDEX idx_device_user_id ON z_device (user_id);

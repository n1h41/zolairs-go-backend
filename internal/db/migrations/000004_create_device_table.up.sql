CREATE TABLE z_device (
  mac_address VARCHAR(17) PRIMARY KEY NOT NULL UNIQUE,
  user_id UUID NOT NULL,
  device_name VARCHAR(255) NOT NULL,
  category VARCHAR(255),
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES z_users(user_id) ON DELETE CASCADE
);

CREATE INDEX idx_device_user_id ON z_device(user_id);

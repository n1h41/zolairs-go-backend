CREATE TABLE IF NOT EXISTS z_users (
  user_id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  phone VARCHAR(50),
  address JSONB,
  parent_id UUID,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (parent_id) REFERENCES z_users(user_id)
);

CREATE INDEX idx_users_email ON z_users(email);
CREATE INDEX idx_users_parent_id ON z_users(parent_id);


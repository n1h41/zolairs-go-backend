ALTER TABLE z_users ADD COLUMN role ENUM('admin', 'user') NOT NULL DEFAULT 'user';

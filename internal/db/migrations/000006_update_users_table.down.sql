ALTER TABLE z_users
    ADD COLUMN IF NOT EXISTS referral_mail varchar(255) DEFAULT NULL;

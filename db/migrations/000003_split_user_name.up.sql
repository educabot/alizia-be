-- 000003_split_user_name.up.sql
-- Split users.name into first_name + last_name

ALTER TABLE users RENAME COLUMN name TO first_name;
ALTER TABLE users ADD COLUMN last_name VARCHAR(255) NOT NULL DEFAULT '';

-- 000003_split_user_name.down.sql
-- Merge first_name + last_name back into name

UPDATE users SET first_name = first_name || ' ' || last_name WHERE last_name != '';
ALTER TABLE users DROP COLUMN last_name;
ALTER TABLE users RENAME COLUMN first_name TO name;

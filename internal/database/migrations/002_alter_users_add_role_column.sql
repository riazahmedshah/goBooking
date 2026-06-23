-- Write your migrate up statements here

ALTER TABLE users ADD COLUMN role TEXT NOT NULL DEFAULT 'user';

ALTER TABLE users ADD CONSTRAINT chk_user_role CHECK (role IN ('user', 'host'));

UPDATE users SET role = 'user' WHERE is_host=false;
UPDATE users SET role = 'host' WHERE is_host=true;

ALTER TABLE users DROP COLUMN is_host;
---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

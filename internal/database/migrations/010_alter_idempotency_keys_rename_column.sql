-- Write your migrate up statements here

ALTER TABLE idempotency_keys RENAME COLUMN is_finalised TO is_finalized;

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

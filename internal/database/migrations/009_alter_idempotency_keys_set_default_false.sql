-- Write your migrate up statements here

ALTER TABLE idempotency_keys ALTER COLUMN is_finalised SET DEFAULT false;

UPDATE idempotency_keys SET is_finalised = false WHERE is_finalised IS NULL;

ALTER TABLE idempotency_keys ALTER COLUMN is_finalised SET NOT NULL;


---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

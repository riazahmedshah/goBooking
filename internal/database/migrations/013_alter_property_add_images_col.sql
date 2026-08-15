-- Write your migrate up statements here

ALTER TABLE properties 
ADD COLUMN images TEXT[] NOT NULL DEFAULT '{}';

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

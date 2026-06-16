-- Write your migrate up statements here

CREATE TABLE properties (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL UNIQUE,
  sub_title TEXT NOT NULL,
  image TEXT NOT NULL,
  address_id INT NOT NULL -- (Later) - REFERENCES adressess(id) ON DELETE RESTRICT
  max_guests INT DEFAULT 1
  ratings INT DEFAULT 0
)

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

-- Write your migrate up statements here

CREATE TABLE properties (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL UNIQUE,
  sub_title TEXT,
  image TEXT,
  address_id INT NOT NULL, -- (Later) - REFERENCES adressess(id) ON DELETE RESTRICT
  host_id INT NOT NULL, -- (Later) - REFERENCES hosts(id) ON DELETE RESTRICT
  max_guests INT DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

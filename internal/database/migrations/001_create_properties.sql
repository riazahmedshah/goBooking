-- Write your migrate up statements here

CREATE TABLE bookings (
  id SERIAL PRIMARY KEY,
  property_id INT NOT NULL,
  user_id INT NOT NULL,
  total_price INT NOT NULL,
  status TEXT
)

CREATE TABLE idempotency_keys (
  id SERIAL PRIMARY KEY,
  idem_key TEXT,
  booking_id INT NOT NULL,
  is_finalized BOOLEAN DEFAULT false
)

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

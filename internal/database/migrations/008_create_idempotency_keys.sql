-- Write your migrate up statements here

CREATE TABLE idempotency_keys (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  booking_id UUID,
  key TEXT UNIQUE NOT NULL,
  is_finalised BOOLEAN,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  CONSTRAINT fk_idepotancy_booking
    FOREIGN KEY (booking_id)
    REFERENCES bookings (id)
    ON DELETE RESTRICT
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

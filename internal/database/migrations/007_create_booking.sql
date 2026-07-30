-- Write your migrate up statements here

CREATE TABLE bookings (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL,
  property_id UUID NOT NULL,
  total_price NUMERIC(10,2) NOT NULL,
  check_in DATE NOT NULL,
  check_out DATE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL, --2026-07-30 13:22:27.96211+00

  status TEXT DEFAULT 'pending' NOT NULL,

  CONSTRAINT chk_status CHECK (status IN ('pending','confirmed','canceled')),

  CONSTRAINT chk_booking_dates CHECK (check_out > check_in),

  CONSTRAINT fk_booking_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE,
  
  CONSTRAINT fk_booking_property
    FOREIGN KEY (property_id)
    REFERENCES properties (id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

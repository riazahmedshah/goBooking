-- Write your migrate up statements here

CREATE TABLE property_availability (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  property_id UUID NOT NULL,
  date DATE NOT NULL,
  is_available BOOLEAN DEFAULT TRUE,
  booking_id UUID NULL,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
  
  CONSTRAINT unique_property_date UNIQUE (property_id, date),

  CONSTRAINT fk_property_availability_property
    FOREIGN KEY (property_id)
    REFERENCES properties (id)
    ON DELETE CASCADE,

  CONSTRAINT fk_property_availability_booking
    FOREIGN KEY (booking_id)
    REFERENCES bookings (id)
    ON DELETE RESTRICT
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

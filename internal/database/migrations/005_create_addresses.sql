-- Write your migrate up statements here

CREATE TABLE addresses (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  country TEXT NOT NULL,
  state TEXT NOT NULL,
  pincode TEXT NOT NULL,
  city TEXT,
  area TEXT NOT NULL,
  property_id UUID NOT NULL UNIQUE,

  CONSTRAINT fk_address_property
    FOREIGN KEY (property_id)
    REFERENCES properties (id)
    ON DELETE CASCADE
);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

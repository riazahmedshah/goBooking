-- Write your migrate up statements here

CREATE TABLE properties (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  title TEXT NOT NULL,
  sub_title TEXT,
  image TEXT,
  address_id UUID UNIQUE NOT NULL,
  host_id UUID NOT NULL,
  max_guests INT NOT NULL DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,

  CONSTRAINT fk_property_address
    FOREIGN KEY (address_id)
    REFERENCES addresses (id)
    ON DELETE RESTRICT, 

  CONSTRAINT fk_property_host
    FOREIGN KEY (host_id)
    REFERENCES hosts (id)
    ON DELETE CASCADE
)

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

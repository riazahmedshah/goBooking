-- Write your migrate up statements here
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  first_name TEXT NOT NULL,
  last_name TEXT,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  is_host BOOLEAN DEFAULT false
);

CREATE TABLE host_profiles (
  id UUID PRIMARY KEY DEFAULT uuidv7(),
  user_id UUID NOT NULL UNIQUE,
  state_name TEXT NOT NULL,
  city TEXT NOT NULL,
  area TEXT NOT NULL,

  CONSTRAINT fk_host_user
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);


---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

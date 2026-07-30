-- Write your migrate up statements here

CREATE OR REPLACE FUNCTION set_default_max_guests()

RETURNS TRIGGER AS $$
BEGIN
  IF NEW.max_guests IS NULL THEN
    NEW.max_guests := 1;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_max_guests_default
BEFORE INSERT ON properties
FOR EACH ROW
EXECUTE FUNCTION set_default_max_guests();

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

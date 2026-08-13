-- Write your migrate up statements here

-- 1. Trigger Function
CREATE OR REPLACE FUNCTION populate_property_availability()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO property_availability (property_id, date, is_available)
    SELECT 
        NEW.id,
        d::date,
        TRUE
    FROM generate_series(
        CURRENT_DATE,
        CURRENT_DATE + INTERVAL '90 days',
        INTERVAL '1 day'
    ) AS d;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Trigger Attachment
CREATE TRIGGER trigger_populate_availability
AFTER INSERT ON properties
FOR EACH ROW
EXECUTE FUNCTION populate_property_availability();

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

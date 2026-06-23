
CREATE TABLE bookings (

  id TEXT PRIMARY KEY,

  property_id INT NOT NULL,

  user_id INT NOT NULL,

  total_price NUMERIC(10, 2) NOT NULL,

  status TEXT DEFAULT 'pending' CHECK (status IN  ('pending', 'confirmed', 'cancelled'))

)



CREATE TABLE idempotency_keys (

  id TEXT PRIMARY KEY,

  idem_key TEXT UNIQUE,

  booking_id INT NOT NULL,

  is_finalized BOOLEAN DEFAULT false

)
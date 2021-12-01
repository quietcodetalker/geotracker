CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username varchar(16) UNIQUE,
    created_at timestamp DEFAULT current_timestamp NOT NULL,
    updated_at timestamp DEFAULT current_timestamp NOT NULL
);

CREATE TABLE positions (
    user_id INT REFERENCES users (id) PRIMARY KEY,
    latitude NUMERIC(11, 8) CHECK (latitude >= -180 AND latitude <= 180) NOT NULL,
    longitude NUMERIC(10, 8) CHECK (longitude >= -90 AND longitude <= 90) NOT NULL,
    created_at timestamp DEFAULT current_timestamp NOT NULL,
    updated_at timestamp DEFAULT current_timestamp NOT NULL
);

CREATE TRIGGER update_updated_at BEFORE UPDATE
    ON users FOR EACH ROW EXECUTE PROCEDURE
    update_updated_at();

CREATE TRIGGER update_updated_at BEFORE UPDATE
    ON positions FOR EACH ROW EXECUTE PROCEDURE
    update_updated_at();

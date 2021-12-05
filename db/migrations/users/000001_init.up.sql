CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE users (
    id SERIAL,
    username varchar(16),
    created_at timestamp DEFAULT current_timestamp NOT NULL,
    updated_at timestamp DEFAULT current_timestamp NOT NULL,

    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_key UNIQUE (username),
    CONSTRAINT users_username_valid CHECK (LENGTH(username) >= 4 AND username ~ '^[a-zA-Z0-9]+$')
);

CREATE TABLE locations (
    user_id INT,
    latitude NUMERIC(11, 8) NOT NULL,
    longitude NUMERIC(10, 8) NOT NULL,
    created_at timestamp DEFAULT current_timestamp NOT NULL,
    updated_at timestamp DEFAULT current_timestamp NOT NULL,

    CONSTRAINT locations_pkey PRIMARY KEY (user_id),
    CONSTRAINT locations_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT locations_latitude_valid CHECK (latitude >= -180 AND latitude <= 180),
    CONSTRAINT locations_longitude_valid CHECK (longitude >= -90 AND longitude <= 90)
);

CREATE TRIGGER update_updated_at BEFORE UPDATE
    ON users FOR EACH ROW EXECUTE PROCEDURE
    update_updated_at();

CREATE TRIGGER update_updated_at BEFORE UPDATE
    ON locations FOR EACH ROW EXECUTE PROCEDURE
    update_updated_at();

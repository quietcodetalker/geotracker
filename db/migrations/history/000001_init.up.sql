CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

CREATE OR REPLACE FUNCTION update_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE OR REPLACE FUNCTION fix_points_precision()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.a = point(TRUNC(NEW.a[0]::numeric, 8), TRUNC(NEW.a[1]::numeric, 8));
    NEW.b = point(TRUNC(NEW.b[0]::numeric, 8), TRUNC(NEW.b[1]::numeric, 8));
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE records (
    id SERIAL,
    user_id INT NULL,
    a POINT,
    b POINT,
    created_at timestamp DEFAULT current_timestamp NOT NULL,
    updated_at timestamp DEFAULT current_timestamp NOT NULL,

    CONSTRAINT records_pkey PRIMARY KEY (id),
    CONSTRAINT records_a_longitude_valid CHECK (a[0] >= -180 AND a[0] <= 180),
    CONSTRAINT records_a_latitude_valid CHECK (a[1] >= -90 AND a[1] <= 90),
    CONSTRAINT records_b_longitude_valid CHECK (b[0] >= -180 AND b[0] <= 180),
    CONSTRAINT records_b_latitude_valid CHECK (b[1] >= -90 AND b[1] <= 90)
);

CREATE TRIGGER update_updated_at BEFORE UPDATE
    ON records FOR EACH ROW EXECUTE PROCEDURE
        update_updated_at();

CREATE TRIGGER fix_points_precision BEFORE INSERT OR UPDATE
    ON records FOR EACH ROW EXECUTE PROCEDURE
        fix_points_precision();
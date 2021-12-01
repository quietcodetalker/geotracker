CREATE TABLE transitions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    range DOUBLE PRECISION NOT NULL
);
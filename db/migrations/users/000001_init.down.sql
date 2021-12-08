DROP TRIGGER IF EXISTS update_updated_at ON locations;
DROP TRIGGER IF EXISTS update_updated_at ON users;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS update_updated_at;
DROP EXTENSION IF EXISTS earthdistance;
DROP EXTENSION IF EXISTS cube;

DROP TRIGGER IF EXISTS update_updated_at ON positions;
DROP TRIGGER IF EXISTS update_updated_at ON users;
DROP TABLE IF EXISTS positions;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS update_updated_at;
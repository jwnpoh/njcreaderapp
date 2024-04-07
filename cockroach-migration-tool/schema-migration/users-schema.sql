DROP TABLE IF EXISTS sessions;

CREATE TABLE sessions (
  sessionID UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  token STRING(26) NOT NULL,
  userID UUID NOT NULL,
  expiry TIMESTAMP NOT NULL,
  hash BYTES NOT NULL
);

-- DROP TABLE IF EXISTS users;

-- CREATE TABLE users (
--   id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
--   email STRING(255) UNIQUE NOT NULL,
--   hash BYTES NOT NULL,
--   role STRING(10) NOT NULL,
--   last_login STRING(20) NOT NULL,
--   display_name STRING(255),
--   class STRING(7) NOT NULL
-- );

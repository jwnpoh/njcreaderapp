DROP TABLE IF EXISTS long;

CREATE TABLE long (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  title STRING(255) UNIQUE NOT NULL,
  url STRING(2083) NOT NULL,
  topic STRING(200) NOT NULL
);

CREATE INDEX on long USING GIN (to_tsvector('english', title || topic));

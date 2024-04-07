-- DROP TABLE IF EXISTS articles;

-- CREATE TABLE articles (
--   id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
--   title STRING UNIQUE,
--   url STRING(2083) NOT NULL,
--   topics STRING(200) NOT NULL,
--   questions STRING(255) NOT NULL,
--   question_display TEXT NOT NULL,
--   published_on INT NOT NULL,
--   must_read BOOLEAN NOT NULL
-- );

-- CREATE INDEX on articles USING GIN (to_tsvector('english', title || topics || questions || question_display));

-- DROP TABLE IF EXISTS topics;

-- CREATE TABLE topics (
--   id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
--   topic STRING(255) NOT NULL,
--   article_id UUID
-- );


-- DROP TABLE IF EXISTS questions;

-- CREATE TABLE questions (
--   id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
--   question STRING(10),
--   article_id UUID
-- );

DROP TABLE IF EXISTS question_list;

CREATE TABLE question_list (
  question STRING(10),
  year STRING(4),
  number STRING(2),
  wording STRING(255)
);

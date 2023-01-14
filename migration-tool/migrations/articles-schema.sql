DROP TABLE IF EXISTS articles;

CREATE TABLE articles (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) UNIQUE NOT NULL,
  url VARCHAR(2083) NOT NULL,
  topics VARCHAR(200) NOT NULL,
  questions VARCHAR(255) NOT NULL,
  question_display TEXT NOT NULL,
  published_on INT NOT NULL,
  must_read BOOLEAN NOT NULL,
  FULLTEXT search (title, topics, question_display )
);


DROP TABLE IF EXISTS topics;

CREATE TABLE topics (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  topic VARCHAR(255) NOT NULL,
  article_id INT
);


DROP TABLE IF EXISTS questions;

CREATE TABLE questions (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  question CHAR(10),
  article_id INT
);

DROP TABLE IF EXISTS question_list;

CREATE TABLE question_list (
  question CHAR(10),
  year CHAR(4),
  number CHAR(2),
  wording VARCHAR(255)
);

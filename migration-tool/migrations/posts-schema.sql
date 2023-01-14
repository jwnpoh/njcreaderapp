DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id INT NOT NULL,
  author VARCHAR(50) NOT NULL,
  likes INT NOT NULL,
  tldr VARCHAR(255) NOT NULL,
  examples TEXT NOT NULL,
  notes TEXT,
  tags VARCHAR(255),
  created_at INT NOT NULL,
  public BOOLEAN NOT NULL,
  article_id INT NOT NULL,
  article_title VARCHAR(255) NOT NULL,
  article_url VARCHAR(2083) NOT NULL,
  FULLTEXT search (tldr, examples, notes, tags, author, article_title)
);

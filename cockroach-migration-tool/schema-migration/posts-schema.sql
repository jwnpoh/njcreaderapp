DROP TABLE IF EXISTS notes;

CREATE TABLE notes (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id UUID NOT NULL,
  author STRING(50) NOT NULL,
  author_class STRING(50) NOT NULL,
  likes INT NOT NULL,
  tldr STRING(255) NOT NULL,
  examples TEXT NOT NULL,
  notes TEXT,
  tags STRING(255),
  created_at INT NOT NULL,
  public BOOLEAN NOT NULL,
  article_id UUID NOT NULL,
  article_title STRING(255) NOT NULL,
  article_url STRING(2083) NOT NULL
);

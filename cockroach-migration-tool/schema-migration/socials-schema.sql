DROP TABLE IF EXISTS likes_list;

CREATE TABLE likes_list (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  post_id UUID NOT NULL,
  liked_by INT NOT NULL
);

DROP TABLE IF EXISTS follows;

CREATE TABLE follows (
  id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  user_id UUID NOT NULL,
  follows INT NOT NULL
);

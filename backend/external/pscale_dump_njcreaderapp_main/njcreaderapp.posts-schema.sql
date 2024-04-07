CREATE TABLE `posts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `author` varchar(50) NOT NULL,
  `author_class` varchar(50) NOT NULL,
  `likes` int NOT NULL,
  `tldr` varchar(255) NOT NULL,
  `examples` text NOT NULL,
  `notes` text,
  `tags` varchar(255) DEFAULT NULL,
  `created_at` int NOT NULL,
  `public` tinyint(1) NOT NULL,
  `article_id` int NOT NULL,
  `article_title` varchar(255) NOT NULL,
  `article_url` varchar(2083) NOT NULL,
  PRIMARY KEY (`id`),
  FULLTEXT KEY `search` (`tldr`,`examples`,`notes`,`tags`,`author`,`article_title`)
) ENGINE=InnoDB AUTO_INCREMENT=86 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `articles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `url` varchar(2083) NOT NULL,
  `topics` varchar(200) NOT NULL,
  `questions` varchar(255) NOT NULL,
  `question_display` text NOT NULL,
  `published_on` int NOT NULL,
  `must_read` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `title` (`title`),
  FULLTEXT KEY `search` (`title`,`topics`,`question_display`)
) ENGINE=InnoDB AUTO_INCREMENT=5492 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `questions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `question` char(10) DEFAULT NULL,
  `article_id` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10962 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `topics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `topic` varchar(255) NOT NULL,
  `article_id` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9843 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

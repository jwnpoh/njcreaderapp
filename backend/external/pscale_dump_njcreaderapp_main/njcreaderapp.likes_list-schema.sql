CREATE TABLE `likes_list` (
  `id` int NOT NULL AUTO_INCREMENT,
  `post_id` int NOT NULL,
  `liked_by` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=334 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `hash` varchar(255) NOT NULL,
  `role` varchar(10) NOT NULL,
  `last_login` varchar(20) NOT NULL,
  `display_name` varchar(255) DEFAULT NULL,
  `class` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=2878 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

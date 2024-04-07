CREATE TABLE `sessions` (
  `sessionID` int NOT NULL AUTO_INCREMENT,
  `token` varchar(26) NOT NULL,
  `userID` int NOT NULL,
  `expiry` timestamp NOT NULL,
  `hash` varchar(255) NOT NULL,
  PRIMARY KEY (`sessionID`)
) ENGINE=InnoDB AUTO_INCREMENT=2323 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------
-- Create application database account.
-- --------------------------------------------------------

CREATE USER IF NOT EXISTS 'cmdbd'@'%'
IDENTIFIED BY PASSWORD '*827FDBE6EA30E5FB96D2573FB10EF1B3F49247A2';

REVOKE ALL PRIVILEGES, GRANT OPTION
FROM 'cmdbd'@'%';

GRANT SELECT, INSERT, UPDATE, LOCK TABLES, DELETE, EXECUTE
ON `gocmdb`.*
TO 'cmdbd'@'%';

-- --------------------------------------------------------
-- Create user accounts - PRODUCTION
-- --------------------------------------------------------
/*
USE `gocmdb`;

LOCK TABLES `cmdb_users` WRITE;
INSERT IGNORE INTO `cmdb_users` (
	username,
	password,
	role
)
VALUES (
	'clubpc',
	'$2a$10$H4icOLYN5ZSO.Fqv5RvrkOY2NDMvNN0NtOexgsEQSRDthZrOiIMnK'
	'agent'
);

UNLOCK TABLES;
*/
-- --------------------------------------------------------
-- Create user accounts - NON-PRODUCTION
-- --------------------------------------------------------

USE `gocmdb`;

LOCK TABLES `cmdb_users` WRITE;
INSERT IGNORE INTO `cmdb_users` (
	username,
	password,
	role
)
VALUES (
	'clubpc',
	'$2a$10$Rwh9Ix7Q9.5ST49GngEUJu/VOAYdWG4wnMA9ArSv4qVWQ6nRkyPme',
	'agent'
);

UNLOCK TABLES;

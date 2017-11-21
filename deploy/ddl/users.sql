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
-- Create user accounts.
-- --------------------------------------------------------

USE `gocmdb`;
LOCK TABLES `cmdb_users` WRITE;
INSERT INTO `cmdb_users` (
	username,
	password,
	role
)
VALUES (
	'clubpc',
	'$2a$10$Rwh9Ix7Q9.5ST49GngEUJu/VOAYdWG4wnMA9ArSv4qVWQ6nRkyPme',
	-- '$2a$10$6bSZ98lc/iiZHFKuHyhwJ.IVf.ufyuAfejGWd.QMS721zZtzXfrAC',
	'agent'
);
UNLOCK TABLES;


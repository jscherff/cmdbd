-- --------------------------------------------------------
-- Create application database account.
-- --------------------------------------------------------

CREATE USER IF NOT EXISTS 'cmdbd'@'%';

REVOKE ALL PRIVILEGES, GRANT OPTION
FROM 'cmdbd'@'%';

GRANT SELECT, INSERT, UPDATE, LOCK TABLES, DELETE, EXECUTE
ON `gocmdb`.*
TO 'cmdbd'@'%';

-- --------------------------------------------------------
-- Create application user accounts.
-- --------------------------------------------------------

USE `gocmdb`;

LOCK TABLES `cmdb_users` WRITE;
INSERT IGNORE INTO `cmdb_users` (
	username,
	role
)
VALUES (
	'clubpc',
	'agent'
);

UNLOCK TABLES;

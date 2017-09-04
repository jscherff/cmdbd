-- --------------------------------------------------------
-- Host:                         sysadm-dev-01
-- Server version:               10.2.8-MariaDB - MariaDB Server
-- Server OS:                    Linux
-- HeidiSQL Version:             9.4.0.5125
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- Dumping database structure for gocmdb
DROP DATABASE IF EXISTS `gocmdb`;
CREATE DATABASE IF NOT EXISTS `gocmdb` /*!40100 DEFAULT CHARACTER SET latin1 */;
USE `gocmdb`;

-- Dumping structure for table gocmdb.audits
DROP TABLE IF EXISTS `audits`;
CREATE TABLE IF NOT EXISTS `audits` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `serial_number` varchar(50) NOT NULL,
  `field_name` varchar(50) NOT NULL,
  `old_value` varchar(255) DEFAULT NULL,
  `new_value` varchar(255) DEFAULT NULL,
  `audit_date` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `serial_number` (`serial_number`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.checkins
DROP TABLE IF EXISTS `checkins`;
CREATE TABLE IF NOT EXISTS `checkins` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_number` varchar(126) NOT NULL,
  `vendor_name` varchar(126) DEFAULT NULL,
  `product_name` varchar(126) DEFAULT NULL,
  `product_ver` varchar(7) DEFAULT NULL,
  `software_id` varchar(11) DEFAULT NULL,
  `object_type` varchar(255) NOT NULL DEFAULT 'Generic',
  `checkin_date` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `host_name` (`host_name`),
  KEY `serial_number` (`serial_number`),
  CONSTRAINT `FK_checkins_devices` FOREIGN KEY (`serial_number`) REFERENCES `devices` (`serial_number`) ON DELETE NO ACTION ON UPDATE NO ACTION
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.devices
DROP TABLE IF EXISTS `devices`;
CREATE TABLE IF NOT EXISTS `devices` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_number` varchar(126) NOT NULL,
  `vendor_name` varchar(126) DEFAULT NULL,
  `product_name` varchar(126) DEFAULT NULL,
  `product_ver` varchar(7) DEFAULT NULL,
  `software_id` varchar(11) DEFAULT NULL,
  `object_type` varchar(255) NOT NULL DEFAULT 'Generic',
  `first_seen` datetime NOT NULL DEFAULT current_timestamp(),
  `last_seen` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `checkins` int(10) unsigned NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `serial_number` (`serial_number`),
  KEY `host_name` (`host_name`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.serials
DROP TABLE IF EXISTS `serials`;
CREATE TABLE IF NOT EXISTS `serials` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_number` varchar(126) NOT NULL,
  `vendor_name` varchar(126) DEFAULT NULL,
  `product_name` varchar(126) DEFAULT NULL,
  `product_ver` varchar(7) DEFAULT NULL,
  `software_id` varchar(11) DEFAULT NULL,
  `object_type` varchar(255) NOT NULL DEFAULT 'Generic',
  `issue_date` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `serial_number` (`serial_number`),
  KEY `host_name` (`host_name`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.
-- Dumping structure for trigger gocmdb.before_checkins_insert
DROP TRIGGER IF EXISTS `before_checkins_insert`;
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
DELIMITER //
CREATE TRIGGER `before_checkins_insert` BEFORE INSERT ON `checkins` FOR EACH ROW INSERT INTO devices (
                host_name,
                vendor_id,
                product_id,
                serial_number,
                vendor_name,
                product_name,
                product_ver,
                software_id,
                object_type,
                first_seen,
                last_seen
        )
        VALUES (
                NEW.host_name,
                NEW.vendor_id,
                NEW.product_id,
                NEW.serial_number,
                NULLIF(NEW.vendor_name, ''),
                NULLIF(NEW.product_name, ''),
                NULLIF(NEW.product_ver, ''),
                NULLIF(NEW.software_id, ''),
                NEW.object_type,
                NEW.checkin_date,
                NEW.checkin_date
        )
        ON DUPLICATE KEY UPDATE
                host_name = NEW.host_name,
                vendor_id = NEW.vendor_id,
                product_id = NEW.product_id,
                vendor_name = NULLIF(NEW.vendor_name, ''),
                product_name = NULLIF(NEW.product_name, ''),
                product_ver = NULLIF(NEW.product_ver, ''),
                software_id = NULLIF(NEW.software_id, ''),
                object_type = NEW.object_type,
                last_seen = NEW.checkin_date,
                checkins = checkins + 1//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;

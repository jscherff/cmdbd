-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               10.2.8-MariaDB - mariadb.org binary distribution
-- Server OS:                    Win64
-- HeidiSQL Version:             9.4.0.5125
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- truncate all tables
TRUNCATE TABLE `cmdb_sequence`;
TRUNCATE TABLE `usbci_changes`;
TRUNCATE TABLE `usbci_checkins`;
TRUNCATE TABLE `usbci_serialized`;
TRUNCATE TABLE `usbci_snrequests`;
TRUNCATE TABLE `usbci_unserialized`;
TRUNCATE TABLE `usbmeta_vendor`;
TRUNCATE TABLE `usbmeta_product`;
TRUNCATE TABLE `usbmeta_class`;
TRUNCATE TABLE `usbmeta_subclass`;
TRUNCATE TABLE `usbmeta_protocol`;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;

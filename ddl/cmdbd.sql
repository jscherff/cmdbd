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


-- Dumping database structure for gocmdb
DROP DATABASE IF EXISTS `gocmdb`;
CREATE DATABASE IF NOT EXISTS `gocmdb` /*!40100 DEFAULT CHARACTER SET latin1 */;
USE `gocmdb`;

-- Dumping structure for function gocmdb.func_insert_usbci_snrequests
DROP FUNCTION IF EXISTS `func_insert_usbci_snrequests`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` FUNCTION `func_insert_usbci_snrequests`(
	`host_name_in` varchar(255),
	`vendor_id_in` varchar(4),
	`product_id_in` varchar(4),
	`serial_num_in` varchar(126),
	`vendor_name_in` varchar(126),
	`product_name_in` varchar(126),
	`product_ver_in` varchar(7),
	`software_id_in` varchar(11),
	`bus_number_in` int(10),
	`bus_address_in` int(10),
	`port_number_in` int(10),
	`buffer_size_in` int(10),
	`max_pkt_size_in` int(10),
	`usb_spec_in` varchar(5),
	`usb_class_in` varchar(126),
	`usb_subclass_in` varchar(126),
	`usb_protocol_in` varchar(126),
	`device_speed_in` varchar(126),
	`device_ver_in` varchar(5),
	`device_sn_in` varchar(126),
	`factory_sn_in` varchar(126),
	`descriptor_sn_in` varchar(126),
	`object_type_in` varchar(255),
	`checkin_date_in` DATETIME

) RETURNS int(11)
    DETERMINISTIC
    SQL SECURITY INVOKER
BEGIN
	INSERT INTO usbci_snrequests (
		host_name,
		vendor_id,
		product_id,
		serial_num,
		vendor_name,
		product_name,
		product_ver,
		software_id,
		port_number,
		bus_number,
		bus_address,
		buffer_size,
		max_pkt_size,
		usb_spec,
		usb_class,
		usb_subclass,
		usb_protocol,
		device_speed,
		device_ver,
		device_sn,
		factory_sn,
		descriptor_sn,
		object_type
	)
	VALUES (
		host_name_in,
		vendor_id_in,
		product_id_in,
		serial_num_in,
		vendor_name_in,
		product_name_in,
		product_ver_in,
		software_id_in,
		port_number_in,
		bus_number_in,
		bus_address_in,
		buffer_size_in,
		max_pkt_size_in,
		usb_spec_in,
		usb_class_in,
		usb_subclass_in,
		usb_protocol_in, 
		device_speed_in,
		device_ver_in,
		device_sn_in,
		factory_sn_in,
		descriptor_sn_in,
		object_type_in
	);
	
	RETURN LAST_INSERT_ID();
END//
DELIMITER ;

-- Dumping structure for procedure gocmdb.proc_columns_for_table
DROP PROCEDURE IF EXISTS `proc_columns_for_table`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` PROCEDURE `proc_columns_for_table`(
	IN `table_name_in` VARCHAR(64)


)
    DETERMINISTIC
BEGIN
	SELECT column_name
	FROM information_schema.columns
	WHERE table_name = table_name_in
	AND table_schema = 'gocmdb';
END//
DELIMITER ;

-- Dumping structure for procedure gocmdb.proc_insert_usbci_checkins
DROP PROCEDURE IF EXISTS `proc_insert_usbci_checkins`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` PROCEDURE `proc_insert_usbci_checkins`(
	IN `host_name_in` varchar(255),
	IN `vendor_id_in` varchar(4),
	IN `product_id_in` varchar(4),
	IN `serial_num_in` varchar(126),
	IN `vendor_name_in` varchar(126),
	IN `product_name_in` varchar(126),
	IN `product_ver_in` varchar(7),
	IN `software_id_in` varchar(11),
	IN `bus_number_in` int(10),
	IN `bus_address_in` int(10),
	IN `port_number_in` int(10),
	IN `buffer_size_in` int(10),
	IN `max_pkt_size_in` int(10),
	IN `usb_spec_in` varchar(5),
	IN `usb_class_in` varchar(126),
	IN `usb_subclass_in` varchar(126),
	IN `usb_protocol_in` varchar(126),
	IN `device_speed_in` varchar(126),
	IN `device_ver_in` varchar(5),
	IN `device_sn_in` varchar(126),
	IN `factory_sn_in` varchar(126),
	IN `descriptor_sn_in` varchar(126),
	IN `object_type_in` varchar(255),
	IN `checkin_date_in` DATETIME

)
    DETERMINISTIC
    SQL SECURITY INVOKER
BEGIN
	INSERT INTO usbci_checkins (
		host_name,
		vendor_id,
		product_id,
		serial_num,
		vendor_name,
		product_name,
		product_ver,
		software_id,
		port_number,
		bus_number,
		bus_address,
		buffer_size,
		max_pkt_size,
		usb_spec,
		usb_class,
		usb_subclass,
		usb_protocol,
		device_speed,
		device_ver,
		device_sn,
		factory_sn,
		descriptor_sn,
		object_type
	)
	VALUES (
		host_name_in,
		vendor_id_in,
		product_id_in,
		serial_num_in,
		vendor_name_in,
		product_name_in,
		product_ver_in,
		software_id_in,
		port_number_in,
		bus_number_in,
		bus_address_in,
		buffer_size_in,
		max_pkt_size_in,
		usb_spec_in,
		usb_class_in,
		usb_subclass_in,
		usb_protocol_in, 
		device_speed_in,
		device_ver_in,
		device_sn_in,
		factory_sn_in,
		descriptor_sn_in,
		object_type_in
	);
END//
DELIMITER ;

-- Dumping structure for procedure gocmdb.proc_insert_usbci_serialized
DROP PROCEDURE IF EXISTS `proc_insert_usbci_serialized`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` PROCEDURE `proc_insert_usbci_serialized`(
	IN `host_name_in` varchar(255),
	IN `vendor_id_in` varchar(4),
	IN `product_id_in` varchar(4),
	IN `serial_num_in` varchar(126),
	IN `vendor_name_in` varchar(126),
	IN `product_name_in` varchar(126),
	IN `product_ver_in` varchar(7),
	IN `software_id_in` varchar(11),
	IN `bus_number_in` int(10),
	IN `bus_address_in` int(10),
	IN `port_number_in` int(10),
	IN `buffer_size_in` int(10),
	IN `max_pkt_size_in` int(10),
	IN `usb_spec_in` varchar(5),
	IN `usb_class_in` varchar(126),
	IN `usb_subclass_in` varchar(126),
	IN `usb_protocol_in` varchar(126),
	IN `device_speed_in` varchar(126),
	IN `device_ver_in` varchar(5),
	IN `device_sn_in` varchar(126),
	IN `factory_sn_in` varchar(126),
	IN `descriptor_sn_in` varchar(126),
	IN `object_type_in` varchar(255),
	IN `checkin_date_in` DATETIME

)
    DETERMINISTIC
    SQL SECURITY INVOKER
BEGIN
	INSERT INTO usbci_serialized (
		host_name,
		vendor_id,
		product_id,
		serial_num,
		vendor_name,
		product_name,
		product_ver,
		software_id,
		port_number,
		bus_number,
		bus_address,
		buffer_size,
		max_pkt_size,
		usb_spec,
		usb_class,
		usb_subclass,
		usb_protocol,
		device_speed,
		device_ver,
		device_sn,
		factory_sn,
		descriptor_sn,
		object_type,
		first_seen,
		last_seen
	)
	VALUES (
		host_name_in,
		vendor_id_in,
		product_id_in,
		serial_num_in,
		vendor_name_in,
		product_name_in,
		product_ver_in,
		software_id_in,
		port_number_in,
		bus_number_in,
		bus_address_in,
		buffer_size_in,
		max_pkt_size_in,
		usb_spec_in,
		usb_class_in,
		usb_subclass_in,
		usb_protocol_in, 
		device_speed_in,
		device_ver_in,
		device_sn_in,
		factory_sn_in,
		descriptor_sn_in,
		object_type_in,
		checkin_date_in,
		checkin_date_in
	)
	ON DUPLICATE KEY UPDATE
		-- host_name = host_name_in,
		-- vendor_id = vendor_id_in,
		-- product_id = product_id_in,
		vendor_name = vendor_name_in,
		product_name = product_name_in,
		product_ver = product_ver_in,
		software_id = software_id_in,
		port_number = port_number_in,
		bus_number = bus_number_in,
		bus_address = bus_address_in,
		buffer_size = buffer_size_in,
		max_pkt_size = max_pkt_size_in,
		usb_spec = usb_spec_in,
		usb_class = usb_class_in,
		usb_subclass = usb_subclass_in,
		usb_protocol = usb_protocol_in,
		device_speed = device_speed_in,
		device_ver = device_ver_in,
		device_sn = device_sn_in,
		factory_sn = factory_sn_in,
		descriptor_sn = descriptor_sn_in,
		object_type = object_type_in,
		last_seen = checkin_date_in,
		checkins = checkins + 1;
END//
DELIMITER ;

-- Dumping structure for procedure gocmdb.proc_insert_usbci_unserialized
DROP PROCEDURE IF EXISTS `proc_insert_usbci_unserialized`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` PROCEDURE `proc_insert_usbci_unserialized`(
	IN `host_name_in` varchar(255),
	IN `vendor_id_in` varchar(4),
	IN `product_id_in` varchar(4),
	IN `serial_num_in` varchar(126),
	IN `vendor_name_in` varchar(126),
	IN `product_name_in` varchar(126),
	IN `product_ver_in` varchar(7),
	IN `software_id_in` varchar(11),
	IN `bus_number_in` int(10),
	IN `bus_address_in` int(10),
	IN `port_number_in` int(10),
	IN `buffer_size_in` int(10),
	IN `max_pkt_size_in` int(10),
	IN `usb_spec_in` varchar(5),
	IN `usb_class_in` varchar(126),
	IN `usb_subclass_in` varchar(126),
	IN `usb_protocol_in` varchar(126),
	IN `device_speed_in` varchar(126),
	IN `device_ver_in` varchar(5),
	IN `device_sn_in` varchar(126),
	IN `factory_sn_in` varchar(126),
	IN `descriptor_sn_in` varchar(126),
	IN `object_type_in` varchar(255),
	IN `checkin_date_in` DATETIME

)
    DETERMINISTIC
    SQL SECURITY INVOKER
BEGIN
	INSERT INTO usbci_unserialized (
		host_name,
		vendor_id,
		product_id,
		serial_num,
		vendor_name,
		product_name,
		product_ver,
		software_id,
		port_number,
		bus_number,
		bus_address,
		buffer_size,
		max_pkt_size,
		usb_spec,
		usb_class,
		usb_subclass,
		usb_protocol,
		device_speed,
		device_ver,
		device_sn,
		factory_sn,
		descriptor_sn,
		object_type,
		first_seen,
		last_seen
	)
	VALUES (
		host_name_in,
		vendor_id_in,
		product_id_in,
		serial_num_in,
		vendor_name_in,
		product_name_in,
		product_ver_in,
		software_id_in,
		port_number_in,
		bus_number_in,
		bus_address_in,
		buffer_size_in,
		max_pkt_size_in,
		usb_spec_in,
		usb_class_in,
		usb_subclass_in,
		usb_protocol_in, 
		device_speed_in,
		device_ver_in,
		device_sn_in,
		factory_sn_in,
		descriptor_sn_in,
		object_type_in,
		checkin_date_in,
		checkin_date_in
	)
	ON DUPLICATE KEY UPDATE
		-- host_name = host_name_in,
		-- vendor_id = vendor_id_in,
		-- product_id = product_id_in,
		vendor_name = vendor_name_in,
		product_name = product_name_in,
		product_ver = product_ver_in,
		software_id = software_id_in,
		-- port_number = port_number_in,
		-- bus_number = bus_number_in,
		bus_address = bus_address_in,
		buffer_size = buffer_size_in,
		max_pkt_size = max_pkt_size_in,
		usb_spec = usb_spec_in,
		usb_class = usb_class_in,
		usb_subclass = usb_subclass_in,
		usb_protocol = usb_protocol_in,
		device_speed = device_speed_in,
		device_ver = device_ver_in,
		device_sn = device_sn_in,
		factory_sn = factory_sn_in,
		descriptor_sn = descriptor_sn_in,
		object_type = object_type_in,
		last_seen = checkin_date_in,
		checkins = checkins + 1;
END//
DELIMITER ;

-- Dumping structure for table gocmdb.usbci_changes
DROP TABLE IF EXISTS `usbci_changes`;
CREATE TABLE IF NOT EXISTS `usbci_changes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) NOT NULL,
  `changes` varchar(2048) NOT NULL,
  `audit_date` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `host_name` (`host_name`),
  KEY `vendor_id` (`vendor_id`),
  KEY `product_id` (`product_id`),
  KEY `serial_num` (`serial_num`),
  CONSTRAINT `CONSTRAINT_1` CHECK (json_valid(`changes`))
) ENGINE=InnoDB AUTO_INCREMENT=346 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.usbci_checkins
DROP TABLE IF EXISTS `usbci_checkins`;
CREATE TABLE IF NOT EXISTS `usbci_checkins` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) NOT NULL,
  `vendor_name` varchar(126) NOT NULL,
  `product_name` varchar(126) NOT NULL,
  `product_ver` varchar(7) NOT NULL,
  `software_id` varchar(11) NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `buffer_size` int(10) unsigned NOT NULL,
  `max_pkt_size` int(10) unsigned NOT NULL,
  `usb_spec` varchar(5) NOT NULL,
  `usb_class` varchar(126) NOT NULL,
  `usb_subclass` varchar(126) NOT NULL,
  `usb_protocol` varchar(126) NOT NULL,
  `device_speed` varchar(126) NOT NULL,
  `device_ver` varchar(5) NOT NULL,
  `device_sn` varchar(126) NOT NULL,
  `factory_sn` varchar(126) NOT NULL,
  `descriptor_sn` varchar(126) NOT NULL,
  `object_type` varchar(255) NOT NULL,
  `checkin_date` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `host_name` (`host_name`),
  KEY `serial_num` (`serial_num`),
  KEY `vendor_id` (`vendor_id`),
  KEY `product_id` (`product_id`),
  KEY `software_id` (`software_id`),
  KEY `device_sn` (`device_sn`),
  KEY `factory_sn` (`factory_sn`),
  KEY `descriptor_sn` (`descriptor_sn`),
  KEY `product_ver` (`product_ver`)
) ENGINE=InnoDB AUTO_INCREMENT=908 DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.usbci_serialized
DROP TABLE IF EXISTS `usbci_serialized`;
CREATE TABLE IF NOT EXISTS `usbci_serialized` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) NOT NULL,
  `vendor_name` varchar(126) NOT NULL,
  `product_name` varchar(126) NOT NULL,
  `product_ver` varchar(7) NOT NULL,
  `software_id` varchar(11) NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `buffer_size` int(10) unsigned NOT NULL,
  `max_pkt_size` int(10) unsigned NOT NULL,
  `usb_spec` varchar(5) NOT NULL,
  `usb_class` varchar(126) NOT NULL,
  `usb_subclass` varchar(126) NOT NULL,
  `usb_protocol` varchar(126) NOT NULL,
  `device_speed` varchar(126) NOT NULL,
  `device_ver` varchar(5) NOT NULL,
  `device_sn` varchar(126) NOT NULL,
  `factory_sn` varchar(126) NOT NULL,
  `descriptor_sn` varchar(126) NOT NULL,
  `object_type` varchar(255) NOT NULL,
  `first_seen` datetime NOT NULL DEFAULT current_timestamp(),
  `last_seen` datetime NOT NULL DEFAULT current_timestamp(),
  `checkins` int(10) unsigned NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_id` (`vendor_id`,`product_id`,`serial_num`),
  KEY `software_id` (`software_id`),
  KEY `device_sn` (`device_sn`),
  KEY `factory_sn` (`factory_sn`),
  KEY `descriptor_sn` (`descriptor_sn`),
  KEY `product_ver` (`product_ver`),
  KEY `host_name` (`host_name`)
) ENGINE=InnoDB AUTO_INCREMENT=861 DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.usbci_snrequests
DROP TABLE IF EXISTS `usbci_snrequests`;
CREATE TABLE IF NOT EXISTS `usbci_snrequests` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) NOT NULL,
  `vendor_name` varchar(126) NOT NULL,
  `product_name` varchar(126) NOT NULL,
  `product_ver` varchar(7) NOT NULL,
  `software_id` varchar(11) NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `buffer_size` int(10) unsigned NOT NULL,
  `max_pkt_size` int(10) unsigned NOT NULL,
  `usb_spec` varchar(5) NOT NULL,
  `usb_class` varchar(126) NOT NULL,
  `usb_subclass` varchar(126) NOT NULL,
  `usb_protocol` varchar(126) NOT NULL,
  `device_speed` varchar(126) NOT NULL,
  `device_ver` varchar(5) NOT NULL,
  `device_sn` varchar(126) NOT NULL,
  `factory_sn` varchar(126) NOT NULL,
  `descriptor_sn` varchar(126) NOT NULL,
  `object_type` varchar(255) NOT NULL,
  `request_date` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `host_name` (`host_name`),
  KEY `vendor_id` (`vendor_id`),
  KEY `product_id` (`product_id`),
  KEY `software_id` (`software_id`),
  KEY `device_sn` (`device_sn`),
  KEY `factory_sn` (`factory_sn`),
  KEY `descriptor_sn` (`descriptor_sn`),
  KEY `product_ver` (`product_ver`),
  KEY `serial_num` (`serial_num`)
) ENGINE=InnoDB AUTO_INCREMENT=157 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.usbci_unserialized
DROP TABLE IF EXISTS `usbci_unserialized`;
CREATE TABLE IF NOT EXISTS `usbci_unserialized` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) DEFAULT NULL,
  `vendor_name` varchar(126) DEFAULT NULL,
  `product_name` varchar(126) DEFAULT NULL,
  `product_ver` varchar(7) NOT NULL,
  `software_id` varchar(11) NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `buffer_size` int(10) unsigned NOT NULL,
  `max_pkt_size` int(10) unsigned NOT NULL,
  `usb_spec` varchar(5) NOT NULL,
  `usb_class` varchar(126) NOT NULL,
  `usb_subclass` varchar(126) NOT NULL,
  `usb_protocol` varchar(126) NOT NULL,
  `device_speed` varchar(126) NOT NULL,
  `device_ver` varchar(5) NOT NULL,
  `device_sn` varchar(126) NOT NULL,
  `factory_sn` varchar(126) NOT NULL,
  `descriptor_sn` varchar(126) NOT NULL,
  `object_type` varchar(255) NOT NULL,
  `first_seen` datetime NOT NULL DEFAULT current_timestamp(),
  `last_seen` datetime NOT NULL DEFAULT current_timestamp(),
  `checkins` int(10) unsigned NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_id` (`host_name`,`vendor_id`,`product_id`,`port_number`,`bus_number`),
  KEY `software_id` (`software_id`),
  KEY `device_sn` (`device_sn`),
  KEY `factory_sn` (`factory_sn`),
  KEY `descriptor_sn` (`descriptor_sn`),
  KEY `product_ver` (`product_ver`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;

-- Data exporting was unselected.
-- Dumping structure for view gocmdb.view_usbci_changes_object
DROP VIEW IF EXISTS `view_usbci_changes_object`;
-- Creating temporary table to overcome VIEW dependency errors
CREATE TABLE `view_usbci_changes_object` (
	`host_name` VARCHAR(255) NOT NULL COLLATE 'latin1_swedish_ci',
	`vendor_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`serial_num` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`changes` VARCHAR(2048) NOT NULL COLLATE 'latin1_swedish_ci'
) ENGINE=MyISAM;

-- Dumping structure for view gocmdb.view_usbci_checkins_object
DROP VIEW IF EXISTS `view_usbci_checkins_object`;
-- Creating temporary table to overcome VIEW dependency errors
CREATE TABLE `view_usbci_checkins_object` (
	`host_name` VARCHAR(255) NOT NULL COLLATE 'latin1_swedish_ci',
	`vendor_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`serial_num` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`vendor_name` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_name` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_ver` VARCHAR(7) NOT NULL COLLATE 'latin1_swedish_ci',
	`software_id` VARCHAR(11) NOT NULL COLLATE 'latin1_swedish_ci',
	`port_number` INT(10) UNSIGNED NOT NULL,
	`bus_number` INT(10) UNSIGNED NOT NULL,
	`bus_address` INT(10) UNSIGNED NOT NULL,
	`buffer_size` INT(10) UNSIGNED NOT NULL,
	`max_pkt_size` INT(10) UNSIGNED NOT NULL,
	`usb_spec` VARCHAR(5) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_class` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_subclass` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_protocol` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_speed` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_ver` VARCHAR(5) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`factory_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`descriptor_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`object_type` VARCHAR(255) NOT NULL COLLATE 'latin1_swedish_ci'
) ENGINE=MyISAM;

-- Dumping structure for view gocmdb.view_usbci_serialized_json
DROP VIEW IF EXISTS `view_usbci_serialized_json`;
-- Creating temporary table to overcome VIEW dependency errors
CREATE TABLE `view_usbci_serialized_json` (
	`vendor_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`serial_num` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`json_object` TEXT NOT NULL COLLATE 'latin1_swedish_ci'
) ENGINE=MyISAM;

-- Dumping structure for view gocmdb.view_usbci_serialized_object
DROP VIEW IF EXISTS `view_usbci_serialized_object`;
-- Creating temporary table to overcome VIEW dependency errors
CREATE TABLE `view_usbci_serialized_object` (
	`host_name` VARCHAR(255) NOT NULL COLLATE 'latin1_swedish_ci',
	`vendor_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`serial_num` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`vendor_name` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_name` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_ver` VARCHAR(7) NOT NULL COLLATE 'latin1_swedish_ci',
	`software_id` VARCHAR(11) NOT NULL COLLATE 'latin1_swedish_ci',
	`port_number` INT(10) UNSIGNED NOT NULL,
	`bus_number` INT(10) UNSIGNED NOT NULL,
	`bus_address` INT(10) UNSIGNED NOT NULL,
	`buffer_size` INT(10) UNSIGNED NOT NULL,
	`max_pkt_size` INT(10) UNSIGNED NOT NULL,
	`usb_spec` VARCHAR(5) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_class` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_subclass` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_protocol` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_speed` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_ver` VARCHAR(5) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`factory_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`descriptor_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`object_type` VARCHAR(255) NOT NULL COLLATE 'latin1_swedish_ci'
) ENGINE=MyISAM;

-- Dumping structure for view gocmdb.view_usbci_snrequests_object
DROP VIEW IF EXISTS `view_usbci_snrequests_object`;
-- Creating temporary table to overcome VIEW dependency errors
CREATE TABLE `view_usbci_snrequests_object` (
	`host_name` VARCHAR(255) NOT NULL COLLATE 'latin1_swedish_ci',
	`vendor_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_id` VARCHAR(4) NOT NULL COLLATE 'latin1_swedish_ci',
	`serial_num` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`vendor_name` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_name` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`product_ver` VARCHAR(7) NOT NULL COLLATE 'latin1_swedish_ci',
	`software_id` VARCHAR(11) NOT NULL COLLATE 'latin1_swedish_ci',
	`port_number` INT(10) UNSIGNED NOT NULL,
	`bus_number` INT(10) UNSIGNED NOT NULL,
	`bus_address` INT(10) UNSIGNED NOT NULL,
	`buffer_size` INT(10) UNSIGNED NOT NULL,
	`max_pkt_size` INT(10) UNSIGNED NOT NULL,
	`usb_spec` VARCHAR(5) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_class` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_subclass` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`usb_protocol` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_speed` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_ver` VARCHAR(5) NOT NULL COLLATE 'latin1_swedish_ci',
	`device_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`factory_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`descriptor_sn` VARCHAR(126) NOT NULL COLLATE 'latin1_swedish_ci',
	`object_type` VARCHAR(255) NOT NULL COLLATE 'latin1_swedish_ci'
) ENGINE=MyISAM;

-- Dumping structure for trigger gocmdb.trig_insert_after_usbci_checkins
DROP TRIGGER IF EXISTS `trig_insert_after_usbci_checkins`;
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
DELIMITER //
CREATE TRIGGER `trig_insert_after_usbci_checkins` AFTER INSERT ON `usbci_checkins` FOR EACH ROW BEGIN
	IF (NEW.serial_num != '') THEN
		CALL proc_insert_usbci_serialized(
			NEW.host_name,
			NEW.vendor_id,
			NEW.product_id,
			NEW.serial_num,
			NEW.vendor_name,
			NEW.product_name,
			NEW.product_ver,
			NEW.software_id,
			NEW.port_number,
			NEW.bus_number,
			NEW.bus_address,
			NEW.buffer_size,
			NEW.max_pkt_size,
			NEW.usb_spec,
			NEW.usb_class,
			NEW.usb_subclass,
			NEW.usb_protocol,
			NEW.device_speed,
			NEW.device_ver,
			NEW.device_sn,
			NEW.factory_sn,
			NEW.descriptor_sn,
			NEW.object_type,
			NEW.checkin_date
		);
	ELSE
		CALL proc_insert_usbci_unserialized(
			NEW.host_name,
			NEW.vendor_id,
			NEW.product_id,
			NEW.serial_num,
			NEW.vendor_name,
			NEW.product_name,
			NEW.product_ver,
			NEW.software_id,
			NEW.port_number,
			NEW.bus_number,
			NEW.bus_address,
			NEW.buffer_size,
			NEW.max_pkt_size,
			NEW.usb_spec,
			NEW.usb_class,
			NEW.usb_subclass,
			NEW.usb_protocol,
			NEW.device_speed,
			NEW.device_ver,
			NEW.device_sn,
			NEW.factory_sn,
			NEW.descriptor_sn,
			NEW.object_type,
			NEW.checkin_date
		);
	END IF;
END//
DELIMITER ;
SET SQL_MODE=@OLDTMP_SQL_MODE;

-- Dumping structure for view gocmdb.view_usbci_changes_object
DROP VIEW IF EXISTS `view_usbci_changes_object`;
-- Removing temporary table and create final VIEW structure
DROP TABLE IF EXISTS `view_usbci_changes_object`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` VIEW `view_usbci_changes_object` AS SELECT
	host_name,
	vendor_id,
	product_id,
	serial_num,
	changes
FROM
	usbci_changes ;

-- Dumping structure for view gocmdb.view_usbci_checkins_object
DROP VIEW IF EXISTS `view_usbci_checkins_object`;
-- Removing temporary table and create final VIEW structure
DROP TABLE IF EXISTS `view_usbci_checkins_object`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` VIEW `view_usbci_checkins_object` AS SELECT
	host_name,
	vendor_id,
	product_id,
	serial_num,
	vendor_name,
	product_name,
	product_ver,
	software_id,
	port_number,
	bus_number,
	bus_address,
	buffer_size,
	max_pkt_size,
	usb_spec,
	usb_class,
	usb_subclass,
	usb_protocol,
	device_speed,
	device_ver,
	device_sn,
	factory_sn,
	descriptor_sn,
	object_type
FROM
	usbci_checkins ;

-- Dumping structure for view gocmdb.view_usbci_serialized_json
DROP VIEW IF EXISTS `view_usbci_serialized_json`;
-- Removing temporary table and create final VIEW structure
DROP TABLE IF EXISTS `view_usbci_serialized_json`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` VIEW `view_usbci_serialized_json` AS SELECT
	vendor_id,
	product_id,
	serial_num,
	JSON_OBJECT(
		'host_name',     host_name,
		'vendor_id',     vendor_id,
		'product_id',    product_id,
		'serial_num',    serial_num,
		'vendor_name',   vendor_name,
		'product_name',  product_name,
		'product_ver',   product_ver,
		'software_id',   software_id,
		'port_number',   port_number,
		'bus_number',    bus_number,
		'bus_address',   bus_address,
		'buffer_size',   buffer_size,
		'max_pkt_size',  max_pkt_size,
		'usb_spec',      usb_spec,
		'usb_class',     usb_class,
		'usb_subclass',  usb_subclass,
		'usb_protocol',  usb_protocol,
		'device_speed',  device_speed,
		'device_ver',    device_ver,
		'device_sn',     device_sn,
		'factory_sn',    factory_sn,
		'descriptor_sn', descriptor_sn,
		'object_type',   object_type
	)
AS
		'json_object'
FROM
	usbci_serialized ;

-- Dumping structure for view gocmdb.view_usbci_serialized_object
DROP VIEW IF EXISTS `view_usbci_serialized_object`;
-- Removing temporary table and create final VIEW structure
DROP TABLE IF EXISTS `view_usbci_serialized_object`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` VIEW `view_usbci_serialized_object` AS SELECT
	host_name,
	vendor_id,
	product_id,
	serial_num,
	vendor_name,
	product_name,
	product_ver,
	software_id,
	port_number,
	bus_number,
	bus_address,
	buffer_size,
	max_pkt_size,
	usb_spec,
	usb_class,
	usb_subclass,
	usb_protocol,
	device_speed,
	device_ver,
	device_sn,
	factory_sn,
	descriptor_sn,
	object_type
FROM
	usbci_serialized ;

-- Dumping structure for view gocmdb.view_usbci_snrequests_object
DROP VIEW IF EXISTS `view_usbci_snrequests_object`;
-- Removing temporary table and create final VIEW structure
DROP TABLE IF EXISTS `view_usbci_snrequests_object`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` VIEW `view_usbci_snrequests_object` AS SELECT
	host_name,
	vendor_id,
	product_id,
	serial_num,
	vendor_name,
	product_name,
	product_ver,
	software_id,
	port_number,
	bus_number,
	bus_address,
	buffer_size,
	max_pkt_size,
	usb_spec,
	usb_class,
	usb_subclass,
	usb_protocol,
	device_speed,
	device_ver,
	device_sn,
	factory_sn,
	descriptor_sn,
	object_type
FROM
	usbci_snrequests ;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;

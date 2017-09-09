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

-- Dumping structure for procedure gocmdb.cascade_serialized_devices
DROP PROCEDURE IF EXISTS `cascade_serialized_devices`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` PROCEDURE `cascade_serialized_devices`(
	IN `host_name_in` varchar(255),
	IN `vendor_id_in` varchar(4),
	IN `product_id_in` varchar(4),
	IN `serial_num_in` varchar(126),
	IN `vendor_name_in` varchar(126),
	IN `product_name_in` varchar(126),
	IN `product_ver_in` varchar(7),
	IN `software_id_in` varchar(11),
	IN `buffer_size_in` int(10),
	IN `bus_number_in` int(10),
	IN `bus_address_in` int(10),
	IN `port_number_in` int(10),
	IN `usb_spec_in` varchar(5),
	IN `usb_class_in` varchar(126),
	IN `usb_subclass_in` varchar(126),
	IN `usb_protocol_in` varchar(126),
	IN `device_speed_in` varchar(126),
	IN `device_ver_in` varchar(5),
	IN `max_pkt_size_in` int(10),
	IN `device_sn_in` varchar(126),
	IN `factory_sn_in` varchar(126),
	IN `descriptor_sn_in` varchar(126),
	IN `object_type_in` varchar(255)









,
	IN `checkin_date_in` DATETIME




)
    SQL SECURITY INVOKER
BEGIN
        INSERT INTO serialized_devices (
                host_name,
                vendor_id,
                product_id,
                serial_num,
                vendor_name,
                product_name,
                product_ver,
                software_id,
                buffer_size,
                bus_number,
                bus_address,
                port_number,
                usb_spec,
                usb_class,
                usb_subclass,
                usb_protocol,
                device_speed,
                device_ver,
                max_pkt_size,
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
                NULLIF(vendor_name_in, ''),
                NULLIF(product_name_in, ''),
                product_ver_in,
                software_id_in,
                buffer_size_in,
                bus_number_in,
                bus_address_in,
                port_number_in,
                NULLIF(usb_spec_in, ''),
                NULLIF(usb_class_in, ''),
                NULLIF(usb_subclass_in, ''),
                NULLIF(usb_protocol_in, ''),
                NULLIF(device_speed_in, ''),
                NULLIF(device_ver_in, ''),
                max_pkt_size_in,
                device_sn_in,
                factory_sn_in,
                descriptor_sn_in,
                NULLIF(object_type_in, ''),
                checkin_date_in,
                checkin_date_in
        )
        ON DUPLICATE KEY UPDATE
                host_name = host_name_in,
                -- vendor_id = vendor_id_in,
                -- product_id = product_id_in,
                -- serial_num = serial_num_in,
                vendor_name = NULLIF(vendor_name_in, ''),
                product_name = NULLIF(product_name_in, ''),
                product_ver = product_ver_in,
                software_id = software_id_in,
                buffer_size = buffer_size_in,
                bus_number = bus_number_in,
                bus_address = bus_address_in,
                port_number = port_number_in,
                usb_spec = NULLIF(usb_spec_in, ''),
                usb_class = NULLIF(usb_class_in, ''),
                usb_subclass = NULLIF(usb_subclass_in, ''),
                usb_protocol = NULLIF(usb_protocol_in, ''),
                device_speed = NULLIF(device_speed_in, ''),
                device_ver = NULLIF(device_ver_in, ''),
                max_pkt_size = max_pkt_size_in,
                device_sn = device_sn_in,
                factory_sn = factory_sn_in,
                descriptor_sn = descriptor_sn_in,
                object_type = NULLIF(object_type_in, ''),
                last_seen = checkin_date_in,
                checkins = checkins + 1;
END//
DELIMITER ;

-- Dumping structure for procedure gocmdb.cascade_unserialized_devices
DROP PROCEDURE IF EXISTS `cascade_unserialized_devices`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` PROCEDURE `cascade_unserialized_devices`(
	IN `host_name_in` varchar(255),
	IN `vendor_id_in` varchar(4),
	IN `product_id_in` varchar(4),
	IN `serial_num_in` varchar(126),
	IN `vendor_name_in` varchar(126),
	IN `product_name_in` varchar(126),
	IN `product_ver_in` varchar(7),
	IN `software_id_in` varchar(11),
	IN `buffer_size_in` int(10),
	IN `bus_number_in` int(10),
	IN `bus_address_in` int(10),
	IN `port_number_in` int(10),
	IN `usb_spec_in` varchar(5),
	IN `usb_class_in` varchar(126),
	IN `usb_subclass_in` varchar(126),
	IN `usb_protocol_in` varchar(126),
	IN `device_speed_in` varchar(126),
	IN `device_ver_in` varchar(5),
	IN `max_pkt_size_in` int(10),
	IN `device_sn_in` varchar(126),
	IN `factory_sn_in` varchar(126),
	IN `descriptor_sn_in` varchar(126),
	IN `object_type_in` varchar(255)









,
	IN `checkin_date_in` DATETIME







)
    SQL SECURITY INVOKER
BEGIN
        INSERT INTO unserialized_devices (
                host_name,
                vendor_id,
                product_id,
                serial_num,
                vendor_name,
                product_name,
                product_ver,
                software_id,
                buffer_size,
                bus_number,
                bus_address,
                port_number,
                usb_spec,
                usb_class,
                usb_subclass,
                usb_protocol,
                device_speed,
                device_ver,
                max_pkt_size,
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
                NULLIF(serial_num_in, ''),
                NULLIF(vendor_name_in, ''),
                NULLIF(product_name_in, ''),
                product_ver_in,
                software_id_in,
                buffer_size_in,
                bus_number_in,
                bus_address_in,
                port_number_in,
                NULLIF(usb_spec_in, ''),
                NULLIF(usb_class_in, ''),
                NULLIF(usb_subclass_in, ''),
                NULLIF(usb_protocol_in, ''),
                NULLIF(device_speed_in, ''),
                NULLIF(device_ver_in, ''),
                max_pkt_size_in,
                device_sn_in,
                factory_sn_in,
                descriptor_sn_in,
                NULLIF(object_type_in, ''),
                checkin_date_in,
                checkin_date_in
        )
        ON DUPLICATE KEY UPDATE
                -- host_name = host_name_in,
                -- vendor_id = vendor_id_in,
                -- product_id = product_id_in,
                vendor_name = NULLIF(vendor_name_in, ''),
                product_name = NULLIF(product_name_in, ''),
                product_ver = product_ver_in,
                software_id = software_id_in,
                buffer_size = buffer_size_in,
                -- bus_number = bus_number_in,
                -- bus_address = bus_address_in,
                -- port_number = port_number_in,
                usb_spec = NULLIF(usb_spec_in, ''),
                usb_class = NULLIF(usb_class_in, ''),
                usb_subclass = NULLIF(usb_subclass_in, ''),
                usb_protocol = NULLIF(usb_protocol_in, ''),
                device_speed = NULLIF(device_speed_in, ''),
                device_ver = NULLIF(device_ver_in, ''),
                max_pkt_size = max_pkt_size_in,
                device_sn = device_sn_in,
                factory_sn = factory_sn_in,
                descriptor_sn = descriptor_sn_in,
                object_type = NULLIF(object_type_in, ''),
                last_seen = checkin_date_in,
                checkins = checkins + 1;
END//
DELIMITER ;

-- Dumping structure for table gocmdb.device_changes
DROP TABLE IF EXISTS `device_changes`;
CREATE TABLE IF NOT EXISTS `device_changes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `field_name` varchar(64) NOT NULL,
  `old_value` varchar(255) NOT NULL,
  `new_value` varchar(255) NOT NULL,
  `audit_date` datetime NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `serial_num` (`serial_num`),
  KEY `field_name` (`field_name`),
  KEY `old_value` (`old_value`),
  KEY `new_value` (`new_value`),
  KEY `host_name` (`host_name`),
  KEY `vendor_id` (`vendor_id`),
  KEY `product_id` (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.device_checkins
DROP TABLE IF EXISTS `device_checkins`;
CREATE TABLE IF NOT EXISTS `device_checkins` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) NOT NULL,
  `vendor_name` varchar(126) DEFAULT NULL,
  `product_name` varchar(126) DEFAULT NULL,
  `product_ver` varchar(7) NOT NULL,
  `software_id` varchar(11) NOT NULL,
  `buffer_size` int(10) unsigned NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `usb_spec` varchar(5) DEFAULT NULL,
  `usb_class` varchar(126) DEFAULT NULL,
  `usb_subclass` varchar(126) DEFAULT NULL,
  `usb_protocol` varchar(126) DEFAULT NULL,
  `device_speed` varchar(126) DEFAULT NULL,
  `device_ver` varchar(5) DEFAULT NULL,
  `max_pkt_size` int(10) unsigned DEFAULT NULL,
  `device_sn` varchar(126) NOT NULL,
  `factory_sn` varchar(126) NOT NULL,
  `descriptor_sn` varchar(126) NOT NULL,
  `object_type` varchar(255) DEFAULT NULL,
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
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;

-- Data exporting was unselected.
-- Dumping structure for procedure gocmdb.insert_device_change
DROP PROCEDURE IF EXISTS `insert_device_change`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` PROCEDURE `insert_device_change`(
	IN `host_name_in` varchar(255),
	IN `vendor_id_in` varchar(4),
	IN `product_id_in` varchar(4),
	IN `serial_num_in` varchar(126),
	IN `bus_number_in` int(10),
	IN `bus_address_in` int(10),
	IN `port_number_in` int(10),
	IN `field_name_in` varchar(64),
	IN `old_value_in` varchar(255),
	IN `new_value_in` varchar(255)










)
    DETERMINISTIC
    SQL SECURITY INVOKER
BEGIN
        INSERT INTO device_changes (
                host_name,
                vendor_id,
                product_id,
                serial_num,
                bus_number,
                bus_address,
                port_number,
                field_name,
                old_value,
                new_value
        )
        VALUES (
                host_name_in,
                vendor_id_in,
                product_id_in,
                serial_num_in,
                bus_number_in,
                bus_address_in,
                port_number_in,
                field_name_in,
                old_value_in,
                new_value_in
        );
END//
DELIMITER ;

-- Dumping structure for function gocmdb.insert_device_checkin
DROP FUNCTION IF EXISTS `insert_device_checkin`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` FUNCTION `insert_device_checkin`(
	`host_name_in` varchar(255),
	`vendor_id_in` varchar(4),
	`product_id_in` varchar(4),
	`serial_num_in` varchar(126),
	`vendor_name_in` varchar(126),
	`product_name_in` varchar(126),
	`product_ver_in` varchar(7),
	`software_id_in` varchar(11),
	`buffer_size_in` int(10),
	`bus_number_in` int(10),
	`bus_address_in` int(10),
	`port_number_in` int(10),
	`usb_spec_in` varchar(5),
	`usb_class_in` varchar(126),
	`usb_subclass_in` varchar(126),
	`usb_protocol_in` varchar(126),
	`device_speed_in` varchar(126),
	`device_ver_in` varchar(5),
	`max_pkt_size_in` int(10),
	`device_sn_in` varchar(126),
	`factory_sn_in` varchar(126),
	`descriptor_sn_in` varchar(126),
	`object_type_in` varchar(255)










) RETURNS int(11)
    DETERMINISTIC
    SQL SECURITY INVOKER
BEGIN
        INSERT INTO device_checkins (
                host_name,
                vendor_id,
                product_id,
                serial_num,
                vendor_name,
                product_name,
                product_ver,
                software_id,
                buffer_size,
                bus_number,
                bus_address,
                port_number,
                usb_spec,
                usb_class,
                usb_subclass,
                usb_protocol,
                device_speed,
                device_ver,
                max_pkt_size,
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
                NULLIF(vendor_name_in, ''),
                NULLIF(product_name_in, ''),
                product_ver_in,
                software_id_in,
                buffer_size_in,
                bus_number_in,
                bus_address_in,
                port_number_in,
                NULLIF(usb_spec_in, ''),
                NULLIF(usb_class_in, ''),
                NULLIF(usb_subclass_in, ''),
                NULLIF(usb_protocol_in, ''),
                NULLIF(device_speed_in, ''),
                NULLIF(device_ver_in, ''),
                max_pkt_size_in,
                device_sn_in,
                factory_sn_in,
                descriptor_sn_in,
                NULLIF(object_type_in, '')
        );
        
        RETURN LAST_INSERT_ID();
END//
DELIMITER ;

-- Dumping structure for function gocmdb.insert_serial_request
DROP FUNCTION IF EXISTS `insert_serial_request`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` FUNCTION `insert_serial_request`(
	`host_name_in` varchar(255),
	`vendor_id_in` varchar(4),
	`product_id_in` varchar(4),
	`serial_num_in` varchar(126),
	`vendor_name_in` varchar(126),
	`product_name_in` varchar(126),
	`product_ver_in` varchar(7),
	`software_id_in` varchar(11),
	`buffer_size_in` int(10),
	`bus_number_in` int(10),
	`bus_address_in` int(10),
	`port_number_in` int(10),
	`usb_spec_in` varchar(5),
	`usb_class_in` varchar(126),
	`usb_subclass_in` varchar(126),
	`usb_protocol_in` varchar(126),
	`device_speed_in` varchar(126),
	`device_ver_in` varchar(5),
	`max_pkt_size_in` int(10),
	`device_sn_in` varchar(126),
	`factory_sn_in` varchar(126),
	`descriptor_sn_in` varchar(126),
	`object_type_in` varchar(255)












) RETURNS int(10)
    DETERMINISTIC
    SQL SECURITY INVOKER
BEGIN
        INSERT INTO serial_requests (
                host_name,
                vendor_id,
                product_id,
                serial_num,
                vendor_name,
                product_name,
                product_ver,
                software_id,
                buffer_size,
                bus_number,
                bus_address,
                port_number,
                usb_spec,
                usb_class,
                usb_subclass,
                usb_protocol,
                device_speed,
                device_ver,
                max_pkt_size,
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
                NULLIF(vendor_name_in, ''),
                NULLIF(product_name_in, ''),
                product_ver_in,
                software_id_in,
                buffer_size_in,
                bus_number_in,
                bus_address_in,
                port_number_in,
                NULLIF(usb_spec_in, ''),
                NULLIF(usb_class_in, ''),
                NULLIF(usb_subclass_in, ''),
                NULLIF(usb_protocol_in, ''),
                NULLIF(device_speed_in, ''),
                NULLIF(device_ver_in, ''),
                max_pkt_size_in,
                device_sn_in,
                factory_sn_in,
                descriptor_sn_in,
                NULLIF(object_type_in, '')
        );
        
        RETURN LAST_INSERT_ID();
END//
DELIMITER ;

-- Dumping structure for table gocmdb.serialized_devices
DROP TABLE IF EXISTS `serialized_devices`;
CREATE TABLE IF NOT EXISTS `serialized_devices` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) NOT NULL,
  `vendor_name` varchar(126) DEFAULT NULL,
  `product_name` varchar(126) DEFAULT NULL,
  `product_ver` varchar(7) NOT NULL,
  `software_id` varchar(11) NOT NULL,
  `buffer_size` int(10) unsigned NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `usb_spec` varchar(5) DEFAULT NULL,
  `usb_class` varchar(126) DEFAULT NULL,
  `usb_subclass` varchar(126) DEFAULT NULL,
  `usb_protocol` varchar(126) DEFAULT NULL,
  `device_speed` varchar(126) DEFAULT NULL,
  `device_ver` varchar(5) DEFAULT NULL,
  `max_pkt_size` int(10) unsigned DEFAULT NULL,
  `device_sn` varchar(126) NOT NULL,
  `factory_sn` varchar(126) NOT NULL,
  `descriptor_sn` varchar(126) NOT NULL,
  `object_type` varchar(255) DEFAULT NULL,
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
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.serial_requests
DROP TABLE IF EXISTS `serial_requests`;
CREATE TABLE IF NOT EXISTS `serial_requests` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) NOT NULL,
  `vendor_name` varchar(126) DEFAULT NULL,
  `product_name` varchar(126) DEFAULT NULL,
  `product_ver` varchar(7) NOT NULL,
  `software_id` varchar(11) NOT NULL,
  `buffer_size` int(10) unsigned NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `usb_spec` varchar(5) DEFAULT NULL,
  `usb_class` varchar(126) DEFAULT NULL,
  `usb_subclass` varchar(126) DEFAULT NULL,
  `usb_protocol` varchar(126) DEFAULT NULL,
  `device_speed` varchar(126) DEFAULT NULL,
  `device_ver` varchar(5) DEFAULT NULL,
  `max_pkt_size` int(10) unsigned DEFAULT NULL,
  `device_sn` varchar(126) NOT NULL,
  `factory_sn` varchar(126) NOT NULL,
  `descriptor_sn` varchar(126) NOT NULL,
  `object_type` varchar(255) DEFAULT NULL,
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
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=latin1;

-- Data exporting was unselected.
-- Dumping structure for table gocmdb.unserialized_devices
DROP TABLE IF EXISTS `unserialized_devices`;
CREATE TABLE IF NOT EXISTS `unserialized_devices` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_name` varchar(255) NOT NULL,
  `vendor_id` varchar(4) NOT NULL,
  `product_id` varchar(4) NOT NULL,
  `serial_num` varchar(126) DEFAULT NULL,
  `vendor_name` varchar(126) DEFAULT NULL,
  `product_name` varchar(126) DEFAULT NULL,
  `product_ver` varchar(7) NOT NULL,
  `software_id` varchar(11) NOT NULL,
  `buffer_size` int(10) unsigned NOT NULL,
  `bus_number` int(10) unsigned NOT NULL,
  `bus_address` int(10) unsigned NOT NULL,
  `port_number` int(10) unsigned NOT NULL,
  `usb_spec` varchar(5) DEFAULT NULL,
  `usb_class` varchar(126) DEFAULT NULL,
  `usb_subclass` varchar(126) DEFAULT NULL,
  `usb_protocol` varchar(126) DEFAULT NULL,
  `device_speed` varchar(126) DEFAULT NULL,
  `device_ver` varchar(5) DEFAULT NULL,
  `max_pkt_size` int(10) unsigned DEFAULT NULL,
  `device_sn` varchar(126) NOT NULL,
  `factory_sn` varchar(126) NOT NULL,
  `descriptor_sn` varchar(126) NOT NULL,
  `object_type` varchar(255) DEFAULT NULL,
  `first_seen` datetime NOT NULL DEFAULT current_timestamp(),
  `last_seen` datetime NOT NULL DEFAULT current_timestamp(),
  `checkins` int(10) unsigned NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_id` (`host_name`,`vendor_id`,`product_id`,`bus_number`,`bus_address`,`port_number`),
  KEY `software_id` (`software_id`),
  KEY `device_sn` (`device_sn`),
  KEY `factory_sn` (`factory_sn`),
  KEY `descriptor_sn` (`descriptor_sn`),
  KEY `product_ver` (`product_ver`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=latin1 ROW_FORMAT=DYNAMIC;

-- Data exporting was unselected.
-- Dumping structure for procedure gocmdb.update_serial_request
DROP PROCEDURE IF EXISTS `update_serial_request`;
DELIMITER //
CREATE DEFINER=`root`@`localhost` PROCEDURE `update_serial_request`(
	IN `serial_num_in` VARCHAR(126),
	IN `id_in` INT(10)





)
    DETERMINISTIC
    SQL SECURITY INVOKER
BEGIN
	UPDATE serial_requests
	SET serial_num = serial_num_in
	WHERE id = id_in;
END//
DELIMITER ;

-- Dumping structure for trigger gocmdb.after_device_checkin_insert
DROP TRIGGER IF EXISTS `after_device_checkin_insert`;
SET @OLDTMP_SQL_MODE=@@SQL_MODE, SQL_MODE='STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION';
DELIMITER //
CREATE TRIGGER `after_device_checkin_insert` AFTER INSERT ON `device_checkins` FOR EACH ROW BEGIN
        IF (NEW.serial_num != '') THEN
                CALL cascade_serialized_devices(
                        NEW.host_name,
                        NEW.vendor_id,
                        NEW.product_id,
                        NEW.serial_num,
                        NEW.vendor_name,
                        NEW.product_name,
                        NEW.product_ver,
                        NEW.software_id,
                        NEW.buffer_size,
                        NEW.bus_number,
                        NEW.bus_address,
                        NEW.port_number,
                        NEW.usb_spec,
                        NEW.usb_class,
                        NEW.usb_subclass,
                        NEW.usb_protocol,
                        NEW.device_speed,
                        NEW.device_ver,
                        NEW.max_pkt_size,
                        NEW.device_sn,
                        NEW.factory_sn,
                        NEW.descriptor_sn,
                        NEW.object_type,
                        NEW.checkin_date
                );
        ELSE
                CALL cascade_unserialized_devices(
                        NEW.host_name,
                        NEW.vendor_id,
                        NEW.product_id,
                        NEW.serial_num,
                        NEW.vendor_name,
                        NEW.product_name,
                        NEW.product_ver,
                        NEW.software_id,
                        NEW.buffer_size,
                        NEW.bus_number,
                        NEW.bus_address,
                        NEW.port_number,
                        NEW.usb_spec,
                        NEW.usb_class,
                        NEW.usb_subclass,
                        NEW.usb_protocol,
                        NEW.device_speed,
                        NEW.device_ver,
                        NEW.max_pkt_size,
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

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;

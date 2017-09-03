CREATE TABLE `audits` (
	`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`serial_number` VARCHAR(50) NOT NULL,
	`field_name` VARCHAR(50) NOT NULL,
	`old_value` VARCHAR(255) NULL DEFAULT NULL,
	`new_value` VARCHAR(255) NULL DEFAULT NULL,
	`audit_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	INDEX `serial_number` (`serial_number`)
)
COLLATE='latin1_swedish_ci'
ENGINE=InnoDB
AUTO_INCREMENT=1
;

CREATE TABLE `checkins` (
	`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`host_name` VARCHAR(255) NOT NULL,
	`vendor_id` VARCHAR(4) NOT NULL,
	`product_id` VARCHAR(4) NOT NULL,
	`serial_number` VARCHAR(126) NOT NULL,
	`vendor_name` VARCHAR(126) NULL DEFAULT NULL,
	`product_name` VARCHAR(126) NULL DEFAULT NULL,
	`product_ver` VARCHAR(7) NULL DEFAULT NULL,
	`software_id` VARCHAR(11) NULL DEFAULT NULL,
	`object_type` VARCHAR(255) NOT NULL DEFAULT 'Generic',
	`checkin_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	INDEX `host_name` (`host_name`),
	INDEX `serial_number` (`serial_number`),
	CONSTRAINT `FK_checkins_devices` FOREIGN KEY (`serial_number`) REFERENCES `devices` (`serial_number`) ON UPDATE NO ACTION ON DELETE NO ACTION
)
COLLATE='latin1_swedish_ci'
ENGINE=InnoDB
ROW_FORMAT=DYNAMIC
AUTO_INCREMENT=1
;

CREATE TABLE `devices` (
	`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`host_name` VARCHAR(255) NOT NULL,
	`vendor_id` VARCHAR(4) NOT NULL,
	`product_id` VARCHAR(4) NOT NULL,
	`serial_number` VARCHAR(126) NOT NULL,
	`vendor_name` VARCHAR(126) NULL DEFAULT NULL,
	`product_name` VARCHAR(126) NULL DEFAULT NULL,
	`product_ver` VARCHAR(7) NULL DEFAULT NULL,
	`software_id` VARCHAR(11) NULL DEFAULT NULL,
	`object_type` VARCHAR(255) NOT NULL DEFAULT 'Generic',
	`first_seen` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`last_seen` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`checkins` INT(10) UNSIGNED NOT NULL DEFAULT '1',
	PRIMARY KEY (`id`),
	UNIQUE INDEX `serial_number` (`serial_number`),
	INDEX `host_name` (`host_name`)
	-- CONSTRAINT `CHK_serial_number_length` CHECK (LENGTH(`serial_number`) != 0)
)
COLLATE='latin1_swedish_ci'
ENGINE=InnoDB
ROW_FORMAT=DYNAMIC
AUTO_INCREMENT=1
;

CREATE TABLE `serials` (
	`id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	`host_name` VARCHAR(255) NOT NULL,
	`vendor_id` VARCHAR(4) NOT NULL,
	`product_id` VARCHAR(4) NOT NULL,
	`serial_number` VARCHAR(126) NOT NULL,
	`vendor_name` VARCHAR(126) NULL DEFAULT NULL,
	`product_name` VARCHAR(126) NULL DEFAULT NULL,
	`product_ver` VARCHAR(7) NULL DEFAULT NULL,
	`software_id` VARCHAR(11) NULL DEFAULT NULL,
	`object_type` VARCHAR(255) NOT NULL DEFAULT 'Generic',
	`issue_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	INDEX `serial_number` (`serial_number`),
	INDEX `host_name` (`host_name`)
)
COLLATE='latin1_swedish_ci'
ENGINE=InnoDB
AUTO_INCREMENT=1
;

CREATE OR REPLACE DEFINER=`root`@`localhost` TRIGGER `before_checkins_insert`
BEFORE INSERT ON `checkins`
FOR EACH ROW
        INSERT INTO devices (
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
		checkins = checkins + 1;

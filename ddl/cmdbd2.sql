DROP VIEW IF EXISTS `view_usbci_changes`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW `view_usbci_changes` AS
SELECT
	`usbci_changes`.`host_name` AS `host_name`,
	`usbci_changes`.`vendor_id` AS `vendor_id`,
	`usbci_changes`.`product_id` AS `product_id`,
	`usbci_changes`.`serial_number` AS `serial_number`,
	`usbci_changes`.`changes` AS `changes`
FROM
	`usbci_changes`;

--
-- Final view structure for view `view_usbci_checkins`
--

DROP VIEW IF EXISTS `view_usbci_checkins`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW `view_usbci_checkins` AS
SELECT
	`usbci_checkins`.`host_name` AS `host_name`,
	`usbci_checkins`.`vendor_id` AS `vendor_id`,
	`usbci_checkins`.`product_id` AS `product_id`,
	`usbci_checkins`.`serial_number` AS `serial_number`,
	`usbci_checkins`.`vendor_name` AS `vendor_name`,
	`usbci_checkins`.`product_name` AS `product_name`,
	`usbci_checkins`.`product_ver` AS `product_ver`,
	`usbci_checkins`.`firmware_ver` AS `firmware_ver`,
	`usbci_checkins`.`software_id` AS `software_id`,
	`usbci_checkins`.`port_number` AS `port_number`,
	`usbci_checkins`.`bus_number` AS `bus_number`,
	`usbci_checkins`.`bus_address` AS `bus_address`,
	`usbci_checkins`.`buffer_size` AS `buffer_size`,
	`usbci_checkins`.`max_pkt_size` AS `max_pkt_size`,
	`usbci_checkins`.`usb_spec` AS `usb_spec`,
	`usbci_checkins`.`usb_class` AS `usb_class`,
	`usbci_checkins`.`usb_subclass` AS `usb_subclass`,
	`usbci_checkins`.`usb_protocol` AS `usb_protocol`,
	`usbci_checkins`.`device_speed` AS `device_speed`,
	`usbci_checkins`.`device_ver` AS `device_ver`,
	`usbci_checkins`.`device_sn` AS `device_sn`,
	`usbci_checkins`.`factory_sn` AS `factory_sn`,
	`usbci_checkins`.`descriptor_sn` AS `descriptor_sn`,
	`usbci_checkins`.`object_type` AS `object_type`,
	`usbci_checkins`.`object_json` AS `object_json`,
	`usbci_checkins`.`remote_addr` AS `remote_addr`
FROM
	`usbci_checkins`;

--
-- Final view structure for view `view_usbci_devices`
--

DROP VIEW IF EXISTS `view_usbci_devices`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW `view_usbci_devices` AS
SELECT distinct
	`usbci_checkins`.`vendor_id` AS `vendor_id`,
	`usbci_checkins`.`product_id` AS `product_id`,
	`usbci_checkins`.`vendor_name` AS `vendor_name`,
	`usbci_checkins`.`product_name` AS `product_name`
FROM 
	`usbci_checkins`
ORDER BY
	`usbci_checkins`.`vendor_id`,
	`usbci_checkins`.`product_id`;

--
-- Final view structure for view `view_usbci_hosts`
--

DROP VIEW IF EXISTS `view_usbci_hosts`;
CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW `view_usbci_hosts` AS
SELECT distinct
	`usbci_checkins`.`host_name` AS `host_name`
FROM 
	`usbci_checkins`
ORDER BY
	`usbci_checkins`.`host_name`;

--
-- Final view structure for view `view_usbci_snrequests`
--

CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY INVOKER VIEW `view_usbci_snrequests` AS
SELECT
	`usbci_snrequests`.`host_name` AS `host_name`,
	`usbci_snrequests`.`vendor_id` AS `vendor_id`,
	`usbci_snrequests`.`product_id` AS `product_id`,
	`usbci_snrequests`.`serial_number` AS `serial_number`,
	`usbci_snrequests`.`vendor_name` AS `vendor_name`,
	`usbci_snrequests`.`product_name` AS `product_name`,
	`usbci_snrequests`.`product_ver` AS `product_ver`,
	`usbci_snrequests`.`firmware_ver` AS `firmware_ver`,
	`usbci_snrequests`.`software_id` AS `software_id`,
	`usbci_snrequests`.`port_number` AS `port_number`,
	`usbci_snrequests`.`bus_number` AS `bus_number`,
	`usbci_snrequests`.`bus_address` AS `bus_address`,
	`usbci_snrequests`.`buffer_size` AS `buffer_size`,
	`usbci_snrequests`.`max_pkt_size` AS `max_pkt_size`,
	`usbci_snrequests`.`usb_spec` AS `usb_spec`,
	`usbci_snrequests`.`usb_class` AS `usb_class`,
	`usbci_snrequests`.`usb_subclass` AS `usb_subclass`,
	`usbci_snrequests`.`usb_protocol` AS `usb_protocol`,
	`usbci_snrequests`.`device_speed` AS `device_speed`,
	`usbci_snrequests`.`device_ver` AS `device_ver`,
	`usbci_snrequests`.`device_sn` AS `device_sn`,
	`usbci_snrequests`.`factory_sn` AS `factory_sn`,
	`usbci_snrequests`.`descriptor_sn` AS `descriptor_sn`,
	`usbci_snrequests`.`object_type` AS `object_type`,
	`usbci_snrequests`.`object_json` AS `object_json`,
	`usbci_snrequests`.`remote_addr` AS `remote_addr`
FROM
	`usbci_snrequests`;

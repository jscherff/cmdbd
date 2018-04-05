-- -------------------------------------------------------------------
-- Full Card Reader Inventory.
-- -------------------------------------------------------------------
DROP VIEW IF EXISTS view_usbci_cardreaders_all;
CREATE SQL SECURITY INVOKER VIEW view_usbci_cardreaders_all AS 
SELECT
  serial_number,
  vendor_id,
  product_id,
  vendor_name,
  product_name,
  CASE
    WHEN
      vendor_id = '0801' AND
      product_id = '0001'
    THEN
      CASE
      WHEN product_ver = 'V05'
        THEN 'Dynamag MagneSafe'
      ELSE 'SureSwipe'
      END
    ELSE 'NA'
  END AS product_ver,
  firmware_ver,
  host_name,
  SUBSTRING_INDEX(remote_addr, ':', 1) AS ip_address,
  first_seen,
  last_seen,
  checkins
FROM
  usbci_serialized
WHERE
  (vendor_id = '0801' AND product_id = '0001') OR
  (vendor_id = '0acd' AND product_id = '2030')
ORDER BY
  serial_number;

-- -------------------------------------------------------------------
-- New Card Readers.
-- -------------------------------------------------------------------

DROP VIEW IF EXISTS view_usbci_cardreaders_new;
CREATE SQL SECURITY INVOKER VIEW view_usbci_cardreaders_new AS
SELECT
  serial_number,
  vendor_id,
  product_id,
  vendor_name,
  product_name,
  CASE
    WHEN
      vendor_id = '0801' AND
      product_id = '0001'
    THEN
      CASE
      WHEN product_ver = 'V05'
        THEN 'Dynamag MagneSafe'
      ELSE 'SureSwipe'
      END
    ELSE 'NA'
  END AS product_ver,
  firmware_ver,
  host_name,
  SUBSTRING_INDEX(remote_addr, ':', 1) AS ip_address,
  first_seen,
  last_seen,
  checkins
FROM
  usbci_serialized
WHERE 
  checkins = 1 AND
  (
    (vendor_id = '0801' AND product_id = '0001') OR
    (vendor_id = '0acd' AND product_id = '2030')
  )
ORDER BY
  serial_number;

-- -------------------------------------------------------------------
-- Missing Card Readers.
-- -------------------------------------------------------------------
DROP VIEW IF EXISTS view_usbci_cardreaders_missing;
CREATE SQL SECURITY INVOKER VIEW view_usbci_cardreaders_missing AS 
SELECT
  serial_number,
  vendor_id,
  product_id,
  vendor_name,
  product_name,
  CASE
    WHEN
      vendor_id = '0801' AND
      product_id = '0001'
    THEN
      CASE
      WHEN product_ver = 'V05'
        THEN 'Dynamag MagneSafe'
      ELSE 'SureSwipe'
      END
    ELSE 'NA'
  END AS product_ver,
  firmware_ver,
  host_name,
  SUBSTRING_INDEX(remote_addr, ':', 1) AS ip_address,
  first_seen,
  last_seen,
  checkins
FROM
  usbci_serialized
WHERE 
  (TO_DAYS(NOW()) - TO_DAYS(last_seen) > 30) AND
  (
    (vendor_id = '0801' AND product_id = '0001') OR
    (vendor_id = '0acd' AND product_id = '2030')
  )
ORDER BY
  serial_number;
  
-- -------------------------------------------------------------------
-- Card Readers with Recent Changes.
-- -------------------------------------------------------------------

DROP VIEW IF EXISTS view_usbci_cardreader_changes;
CREATE SQL SECURITY INVOKER VIEW view_usbci_cardreader_changes AS 
SELECT
  c.serial_number AS serial_number,
  c.vendor_id AS vendor_id,
  c.product_id AS product_id,
  vendor_name,
  product_name,
  CASE
    WHEN
      s.vendor_id = '0801' AND
      s.product_id = '0001'
    THEN
      CASE
      WHEN product_ver = 'V05'
        THEN 'Dynamag MagneSafe'
      ELSE 'SureSwipe'
      END
    ELSE 'NA'
  END AS product_ver,
  c.host_name,
  SUBSTRING_INDEX(c.remote_addr, ':', 1) AS ip_address,
  firmware_ver,
  property_name,
  previous_value,
  current_value,
  change_date,
  first_seen,
  last_seen,
  checkins
FROM
  usbci_changes c,
  usbci_serialized s
WHERE
  c.serial_number = s.serial_number AND
  NOT property_name = 'DescriptorSN' AND
  NOT previous_value = '' AND
  NOT previous_value = '0' AND
  (
    (c.vendor_id = '0801' AND c.product_id = '0001') OR
    (c.vendor_id = '0acd' AND c.product_id = '2030')
  ) AND
  (TO_DAYS(NOW()) - TO_DAYS(change_date) < 30);

-- -------------------------------------------------------------------
-- Unique Devices
-- -------------------------------------------------------------------

DROP VIEW IF EXISTS view_usbci_unique_devices;
CREATE SQL SECURITY INVOKER VIEW view_usbci_unique_devices AS 
SELECT
  vendor_id,
  product_id,
  vendor_name,
  product_name,
  count(*) AS 'checkins'
FROM
  usbci_checkins
GROUP BY
  vendor_id,
  product_id,
  vendor_name,
  product_name
ORDER BY
  vendor_id,
  product_id;

-- -------------------------------------------------------------------
-- Unique Hosts
-- -------------------------------------------------------------------

DROP VIEW IF EXISTS view_usbci_unique_hosts;
CREATE SQL SECURITY INVOKER VIEW IF NOT EXISTS view_usbci_unique_hosts AS 
SELECT
  host_name,
  count(*) AS 'checkins'
FROM
  usbci_checkins
GROUP BY
  host_name
ORDER BY
  host_name;


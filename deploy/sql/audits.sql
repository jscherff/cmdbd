----------------------------------------------------------------------
-- Show new devices since last check-in.
----------------------------------------------------------------------

SELECT
  host_name,
  serial_number,
  vendor_id,
  product_id,
  vendor_name,
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
  product_name,
  firmware_ver,
  first_seen,
  last_seen,
  checkins
FROM
  usbci_serialized
WHERE
  checkins = 1;

----------------------------------------------------------------------
-- Show devices that appeared after a certain date.
----------------------------------------------------------------------

SELECT
  host_name,
  serial_number,
  vendor_id,
  product_id,
  vendor_name,
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
  product_name,
  firmware_ver,
  first_seen,
  last_seen,
  checkins
FROM
  usbci_serialized
WHERE
  first_seen > '2018-03-28';

----------------------------------------------------------------------
-- Show devices that have not been seen in a month.
----------------------------------------------------------------------

SELECT
  host_name,
  serial_number,
  vendor_id,
  product_id,
  vendor_name,
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
  product_name,
  firmware_ver,
  first_seen,
  last_seen,
  checkins
FROM
  usbci_serialized
WHERE
  DATEDIFF(NOW(), last_seen) > 30;

----------------------------------------------------------------------
-- Show devices that have had property changes in the past month.
----------------------------------------------------------------------

SELECT
  c.host_name AS host_name,
  c.serial_number AS serial_number,
  c.vendor_id AS vendor_id,
  c.product_id AS product_id,
  vendor_name,
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
  product_name,
  firmware_ver,
  property_name,
  previous_value,
  current_value,
  change_date
FROM
  usbci_changes c,
  usbci_serialized s
WHERE
  c.serial_number = s.serial_number AND
  DATEDIFF(NOW(), change_date) < 30 AND
  NOT property_name = 'DescriptorSN' AND
  NOT previous_value = '' AND
  NOT previous_value = '0';

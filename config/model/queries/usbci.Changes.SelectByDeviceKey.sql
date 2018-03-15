SELECT	*
FROM	usbci_changes
WHERE	vendor_id = :vendor_id AND
	product_id = :product_id AND
	serial_number = :serial_number

INSERT INTO usbci_audits (
	vendor_id,
	product_id,
	serial_number,
	host_name,
	remote_addr,
	changes
)
VALUES (
	:vendor_id,
	:product_id,
	:serial_number,
	:host_name,
	:remote_addr,
	:changes
)

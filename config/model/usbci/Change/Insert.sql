INSERT INTO usbci_changes (
	audit_id,
	vendor_id,
	product_id,
	serial_number,
	host_name,
	remote_addr,
	property_name,
	previous_value,
	current_value
)
VALUES (
	:audit_id,
	:vendor_id,
	:product_id,
	:serial_number,
	:host_name,
	:remote_addr,
	:property_name,
	:previous_value,
	:current_value
)

SELECT	*
FROM	usbci_unserialized
WHERE	host_name = :host_name AND
	vendor_id = :vendor_id AND
	product_id = :product_id AND
	port_number = :port_number AND
	bus_number = :bus_number

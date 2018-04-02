SELECT	product_name
FROM	usbmeta_product
WHERE	vendor_id = :vendor_id AND
	product_id = :product_id

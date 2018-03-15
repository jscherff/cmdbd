SELECT 	subclass_desc
FROM 	usbmeta_subclass
WHERE 	class_id = :class_id AND
	subclass_id = :subclass_id

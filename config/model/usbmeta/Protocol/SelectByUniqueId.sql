SELECT	protocol_desc
FROM	usbmeta_protocol
WHERE	class_id = :class_id AND
	subclass_id = :subclass_id AND
 	protocol_id = :protocol_id

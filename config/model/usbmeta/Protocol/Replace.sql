REPLACE INTO usbmeta_protocol (
	class_id,
	subclass_id,
	protocol_id,
	protocol_desc
)
VALUES (
	:class_id,
	:subclass_id,
	:protocol_id,
	:protocol_desc
)

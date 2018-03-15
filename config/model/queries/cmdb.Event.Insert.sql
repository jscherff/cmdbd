INSERT INTO cmdb_events (
	code,
	source,
 	description,
 	host_name,
 	remote_addr
)
VALUES (
	:code,
 	:source,
 	:description,
 	:host_name,
 	:remote_addr
)

// Copyright 2017 John Scherff
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/RackSec/srslog"
	"os"
)

const (
	LogSystem = srslog.LOG_LOCAL6|srslog.LOG_INFO
	LogAccess = srslog.LOG_LOCAL7|srslog.LOG_INFO
	LogError = srslog.LOG_LOCAL7|srslog.LOG_ERR

	LogFileFlags = os.O_APPEND|os.O_CREATE|os.O_WRONLY
	LogFileMode = 0644
	LogDirMode = 0755

	HttpBodySizeLimit int64 = 1048576

	AuditInsertSQL string = `

		INSERT INTO audits (
			serial_num,
			field_name,
			old_value,
			new_value
		)

		VALUES (
			?,
			?,
			NULLIF(?, ''),
			NULLIF(?, '')
		)`

	CheckinInsertSQL string = `

		INSERT INTO checkins (
			host_name,
			vendor_id,
			product_id,
			serial_num,
			vendor_name,
			product_name,
			product_ver,
			software_id,
			buffer_size,
			usb_spec,
			usb_class,
			usb_subclass,
			usb_protocol,
			device_speed,
			device_ver,
			max_pkt_size,
			device_sn,
			factory_sn,
			descriptor_sn,
			object_type
		)

		VALUES (
			?,
			?,
			?,
			?,
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			?,
			NULLIF(?, 0),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, 0),
			?,
			?,
			?,
			NULLIF(?, '')
		)`

	SerialInsertSQL string = `
	
		INSERT INTO serials (
			host_name,
			vendor_id,
			product_id,
			serial_num,
			vendor_name,
			product_name,
			product_ver,
			software_id,
			buffer_size,
			usb_spec,
			usb_class,
			usb_subclass,
			usb_protocol,
			device_speed,
			device_ver,
			max_pkt_size,
			device_sn,
			factory_sn,
			descriptor_sn,
			object_type
		)

		VALUES (
			?,
			?,
			?,
			?,
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			?,
			NULLIF(?, 0),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, ''),
			NULLIF(?, 0),
			?,
			?,
			?,
			NULLIF(?, '')
		)`

	SerialUpdateSQL string = `

		UPDATE serials
		SET serial_num = ?
		WHERE id = ?`
)

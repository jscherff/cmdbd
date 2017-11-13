package main

type PrimaryKeyCol struct {
	Id		int64	`db:"id"`
}

type UsbCiIdentCols struct {
	HostName	string	`db:"host_name"`
	VendorID	string	`db:"vendor_id"`
	ProductID	string	`db:"product_id"`
	SerialNum	string	`db:"serial_number"`
}


type UsbCiBaseCols struct {
	UsbCiIdentCols
	VendorName	string	`db:"vendor_name"`
	ProductName	string	`db:"product_name"`
	ProductVer	string	`db:"product_ver"`
	FirmwareVer	string	`db:"firmware_ver"`
	SoftwareID	string	`db:"software_id"`
	PortNumber	int	`db:"port_number"`
	BusNumber	int	`db:"bus_number"`
	BusAddress	int	`db:"bus_address"`
	BufferSize	int	`db:"buffer_size"`
	MaxPktSize	int	`db:"max_pkt_size"`
	USBSpec		string	`db:"usb_spec"`
	USBClass	string	`db:"usb_class"`
	USBSubClass	string	`db:"usb_subclass"`
	USBProtocol	string	`db:"usb_protocol"`
	DeviceSpeed	string	`db:"device_speed"`
	DeviceVer	string	`db:"device_ver"`
	DeviceSN	string	`db:"device_sn"`
	FactorySN	string	`db:"factory_sn"`
	DescriptorSN	string	`db:"descriptor_sn"`
	ObjectType	string	`db:"object_type"`
	ObjectJSON	[]byte	`db:"object_json"`
	RemoteAddr	string	`db:"remote_addr"`
}

type UsbCiCustomCols struct {
	Custom01	string	`db:"custom_01,omitempty"`
	Custom02	string	`db:"custom_02,omitempty"`
	Custom03	string	`db:"custom_03,omitempty"`
	Custom04	string	`db:"custom_04,omitempty"`
	Custom05	string	`db:"custom_05,omitempty"`
	Custom06	string	`db:"custom_06,omitempty"`
	Custom07	string	`db:"custom_07,omitempty"`
	Custom08	string	`db:"custom_08,omitempty"`
	Custom09	string	`db:"custom_09,omitempty"`
	Custom10	string	`db:"custom_10,omitempty"`
}

type UsbCiSnRequestsVw struct {
	UsbCiIdentCols
	UsbCiBaseCols
}

type UsbCiCheckinsVw struct {
	UsbCiIdentCols
	UsbCiBaseCols
}

type UsbCiSerializedVw struct {
	UsbCiIdentCols
	UsbCiBaseCols
}

type UsbCiChangesVw struct {
	UsbCiIdentCols
	Changes		[]byte	`db:"changes"`
}

type CmdbSequenceTab struct {
	Ord		int64	`db:"ord"`
	IssueDate	string	`db:"issue_date"`
}

type CmdbUsersTab struct {
	PrimaryKeyCol
	Username	string	`db:"username"`
	Password	string	`db:"password"`
	Created		string	`db:"created"`
	Locked		bool	`db:"locked"`
	Role		string	`db:"role"`
}

type UsbCiCheckinsTab struct {
	PrimaryKeyCol
	UsbCiIdentCols
	UsbCiBaseCols
	CheckinDate	string	`db:"checkin_date"`
}

type UsbCiSnRequestsTab struct {
	PrimaryKeyCol
	UsbCiIdentCols
	UsbCiBaseCols
	RequestDate	string	`db:"request_date"`
}

type UsbCiSerializedTab struct {
	PrimaryKeyCol
	UsbCiIdentCols
	UsbCiBaseCols
	FirstSeen	string	`db:"first_seen"`
	LastSeen	string	`db:"last_seen"`
	Checkins	int	`db:"checkins"`
}

type UsbCiUnserializedTab struct {
	PrimaryKeyCol
	UsbCiIdentCols
	UsbCiBaseCols
	FirstSeen	string	`db:"first_seen"`
	LastSeen	string	`db:"last_seen"`
	Checkins	int	`db:"checkins"`
}

type UsbCiChangesTab struct {
	PrimaryKeyCol
	UsbCiIdentCols
	Changes		[]byte	`db:"changes"`
	AuditDate	string	`db:"audit_date"`
}

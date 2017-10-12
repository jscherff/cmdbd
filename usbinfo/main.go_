package main

import (
	`bufio`
	`database/sql`
	`encoding/json`
	`fmt`
	`io`
	`io/ioutil`
	`log`
	`net/http`
	`regexp`

	_ `github.com/go-sql-driver/mysql`
)

const (
	UsbInfoSourceURL = `http://www.linux-usb.org/usb.ids`
	UsbInfoFileMode = 0640

	usbVendorSQL = `
		REPLACE INTO cmdb_meta_usb_vendor (
			vendor_id,
			vendor_name
		)
		VALUES (?, ?)
	`
	usbProductSQL = `
		REPLACE INTO cmdb_meta_usb_product (
			vendor_id,
			product_id,
			product_name
		)
		VALUES (?, ?, ?)
	`
	usbClassSQL = `
		REPLACE INTO cmdb_meta_usb_class (
			class_id,
			class_desc
		)
		VALUES (?, ?)
	`

	usbSubclassSQL = `
		REPLACE INTO cmdb_meta_usb_subclass (
			class_id,
			subclass_id,
			subclass_desc
		)
		VALUES (?, ?, ?)
	`

	usbProtocolSQL = `
		REPLACE INTO cmdb_meta_usb_protocol (
			class_id,
			subclass_id,
			protocol_id,
			protocol_desc
		)
		VALUES (?, ?, ?, ?)
	`
)

var (
	// Vendors stores the vendor and product ID mappings.
	Vendors map[string]*Vendor

	// Classes stores the class, subclass and protocol mappings.
	Classes map[string]*Class

	// vendorRgx extracts the vendor IDs and names from the source data.
	vendorRgx = regexp.MustCompile(`^([0-9A-Fa-f]{4})\s+(.+)$`)

	// productRgx extracts the product IDs and names from the source data.
	productRgx = regexp.MustCompile(`^\t([0-9A-Fa-f]{4})\s+(.+)$`)

	// classRgx extracts the class IDs and descriptions from the source data.
	classRgx = regexp.MustCompile(`^C\s+([0-9A-Fa-f]{2})\s+(.+)$`)

	// subclassRgx extracts the subclass IDs and descriptions from the source data.
	subclassRgx = regexp.MustCompile(`^\t([0-9A-Fa-f]{2})\s+(.+)$`)

	// protocolRgx extracts the protocol IDs and descriptions from the source data.
	protocolRgx = regexp.MustCompile(`^\t\t([0-9A-Fa-f]{2})\s+(.+)$`)
)

// Vendor contains the vendor name and mappings to all the vendor's products.
type Vendor struct {
	Name string
	Product	map[string]string
}

// String returns the name of the vendor.
func (v Vendor) String() string {
	return v.Name
}

// Class contains the name of the class and mappings for each subclass.
type Class struct {
	Name     string
	Subclass map[string]*Subclass
}

// String returns the name of the class.
func (c Class) String() string {
	return c.Name
}

// Subclass contains the name of the subclass and any associated protocols.
type Subclass struct {
	Name     string
	Protocol map[string]string
}

// String returns the name of the Subclass.
func (s Subclass) String() string {
	return s.Name
}

// ParseUsbInfo reads USB vendor, product, class, subclass, and protocol
// information line-by-line from an io.Reader source and populates data
// structures for use by the application.
func ParseUsbInfo(i io.Reader) (map[string]*Vendor, map[string]*Class, error) {

	vendors := make(map[string]*Vendor)
	classes := make(map[string]*Class)
	scanner := bufio.NewScanner(i)

	var (
		vendor *Vendor
		class *Class
		subclass *Subclass
	)

	for scanner.Scan() {

		if matches := productRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			if vendor == nil {
				return nil, nil, fmt.Errorf(`product with no vendor`)
			}
			if vendor.Product == nil {
				vendor.Product = make(map[string]string)
			}

			vendor.Product[matches[1]] = matches[2]
			continue
		}

		if matches := vendorRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			vendor = &Vendor{Name: matches[2]}
			vendors[matches[1]] = vendor
			continue
		}

		if matches := protocolRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			if subclass == nil {
				return nil, nil, fmt.Errorf(`protocol with no subclass`)
			}
			if subclass.Protocol == nil {
				subclass.Protocol = make(map[string]string)
			}

			subclass.Protocol[matches[1]] = matches[2]
			continue
		}

		if matches := subclassRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			if class == nil {
				return nil, nil, fmt.Errorf(`subclass with no class`)
			}
			if class.Subclass == nil {
				class.Subclass = make(map[string]*Subclass)
			}

			subclass = &Subclass{Name: matches[2]}
			class.Subclass[matches[1]] = subclass
			continue
		}

		if matches := classRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			class = &Class{Name: matches[2]}
			classes[matches[1]] = class
		}
	}

	return vendors, classes, scanner.Err()
}

// LoadUsbInfo loads previously-saved USB information from disk.
func LoadUsbInfo(f string, t interface{}) error {

	j, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &t); err != nil {
		return err
	}

	return nil
}

// SaveUsbInfo saves current USB information to disk.
func SaveUsbInfo(f string, t interface{}) error {

	j, err := json.Marshal(t)

	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(f, j, UsbInfoFileMode); err != nil {
		return err
	}

	return nil
}

// RefreshUsbInfo replaces the stored vendor and class mappings with
// data loaded from the source URL.
func RefreshUsbInfo(url string) error {

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	ids, cls, err := ParseUsbInfo(resp.Body)

	if err != nil {
		return err
	}

	Vendors = ids
	Classes = cls

	return nil
}

// RefreshDatabaseUsbInfo updates the USB meta tables in the database.
func RefreshDatabaseUsbInfo(db *sql.DB) error {

	var (
		tx *sql.Tx

		usbVendorStmt, usbProductStmt, usbClassStmt,
		usbSubclassStmt, usbProtocolStmt *sql.Stmt

		err error
	)

	if err = db.Ping(); err != nil {
		return err
	}
	if usbVendorStmt, err = db.Prepare(usbVendorSQL); err != nil {
		return err
	}
	if usbProductStmt, err = db.Prepare(usbProductSQL); err != nil {
		return err
	}
	if usbClassStmt, err = db.Prepare(usbClassSQL); err != nil {
		return err
	}
	if usbSubclassStmt, err = db.Prepare(usbSubclassSQL); err != nil {
		return err
	}
	if usbProtocolStmt, err = db.Prepare(usbProtocolSQL); err != nil {
		return err
	}
	if tx, err = db.Begin(); err != nil {
		return err
	}

	DeviceLoop:
	for vid := range Vendors {
		vdesc := Vendors[vid].String()
		if _, err = tx.Stmt(usbVendorStmt).Exec(vid, vdesc); err != nil {
			break DeviceLoop
		}
		for pid := range Vendors[vid].Product {
			pdesc := Vendors[vid].Product[pid]
			if _, err = tx.Stmt(usbProductStmt).Exec(vid, pid, pdesc); err != nil {
				break DeviceLoop
			}
		}
	}

	ClassLoop:
	for cid := range Classes {
		cdesc := Classes[cid].String()
		if _, err = tx.Stmt(usbClassStmt).Exec(cid, cdesc); err != nil {
			break ClassLoop
		}
		for sid := range Classes[cid].Subclass {
			sdesc := Classes[cid].Subclass[sid].String()
			if _, err = tx.Stmt(usbSubclassStmt).Exec(cid, sid, sdesc); err != nil {
				break ClassLoop
			}
			for pid := range Classes[cid].Subclass[sid].Protocol {
				pdesc := Classes[cid].Subclass[sid].Protocol[pid]
				if _, err = tx.Stmt(usbProtocolStmt).Exec(cid, sid, pid, pdesc); err != nil {
					break ClassLoop
				}
			}
		}
	}

	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	} else {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return err
}

func main() {

	if err := RefreshUsbInfo(UsbInfoSourceURL); err != nil {
		log.Fatal(err)
	}

	for vid := range Vendors {
		fmt.Printf("%s\t%s\n", vid, Vendors[vid])
		for pid := range Vendors[vid].Product {
			fmt.Printf("\t%s\t%s\n", pid, Vendors[vid].Product[pid])
		}
	}

	for cid := range Classes {
		fmt.Printf("%s\t%s\n", cid, Classes[cid])
		for sid := range Classes[cid].Subclass {
			fmt.Printf("\t%s\t%s\n", sid, Classes[cid].Subclass[sid])
			for pid := range Classes[cid].Subclass[sid].Protocol {
				fmt.Printf("\t\t%s\t%s\n", pid, Classes[cid].Subclass[sid].Protocol[pid])
			}
		}
	}

	db, err := sql.Open("mysql", "cmdbd:K2Cvg3NeyR@/gocmdb")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := RefreshDatabaseUsbInfo(db); err != nil {
		log.Fatal(err)
	}

	SaveUsbInfo(`vendors.json`, Vendors)
	SaveUsbInfo(`classes.json`, Classes)
}

package main

import (
	`bufio`
	`database/sql`
	`fmt`
	`net/http`
	`log`
	`regexp`

	_ `github.com/go-sql-driver/mysql`
)


const (
	usbVendorSQL = `
		INSERT INTO cmdb_meta_usb_vendor (
			vendor_id,
			vendor_name
		)
		VALUES (?, ?)
	`
	usbProductSQL = `
		INSERT INTO cmdb_meta_usb_product (
			vendor_id,
			product_id,
			product_name
		)
		VALUES (?, ?, ?)
	`
	usbClassSQL = `
		INSERT INTO cmdb_meta_usb_class (
			class_id,
			class_desc
		)
		VALUES (?, ?)
	`

	usbSubclassSQL = `
		INSERT INTO cmdb_meta_usb_subclass (
			class_id,
			subclass_id,
			subclass_desc
		)
		VALUES (?, ?, ?)
	`

	usbProtocolSQL = `
		INSERT INTO cmdb_meta_usb_protocol (
			class_id,
			subclass_id,
			protocol_id,
			protocol_desc
		)
		VALUES (?, ?, ?, ?)
	`
)

var (
	vidRgx = regexp.MustCompile(`^([0-9A-Fa-f]{4})\s{2}(.+)$`)
	pidRgx = regexp.MustCompile(`^\t([0-9A-Fa-f]{4})\s{2}(.+)$`)

	classRgx = regexp.MustCompile(`^C\s([0-9A-Fa-f]{2})\s{2}(.+)$`)
	subclassRgx = regexp.MustCompile(`^\t([0-9A-Fa-f]{2})\s{2}(.+)$`)
	protocolRgx = regexp.MustCompile(`^\t\t([0-9A-Fa-f]{2})\s{2}(.+)$`)

	db *sql.DB
	usbVendorStmt, usbProductStmt, usbClassStmt, usbSubclassStmt, usbProtocolStmt *sql.Stmt
)

type Database struct {
	*sql.DB
	Driver string
}

type Vendor struct {
	Name string
	Product	map[string]*Product
}

func (v Vendor) String() string {
	return v.Name
}

type Product struct {
	Name string
}

func (p Product) String() string {
	return p.Name
}

type Class struct {
	Name string
	Subclass map[string]*Subclass
}

func (c Class) String() string {
	return c.Name
}

type Subclass struct {
	Name string
	Protocol map[string]*Protocol
}

func (s Subclass) String() string {
	return s.Name
}

type Protocol struct {
	Name string
}

func (p Protocol) String() string {
	return p.Name
}

func init() {

	var err error

	if db, err = sql.Open("mysql", "cmdbd:K2Cvg3NeyR@/gocmdb"); err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	if usbVendorStmt, err = db.Prepare(usbVendorSQL); err != nil {
		log.Fatal(err)
	}
	if usbProductStmt, err = db.Prepare(usbProductSQL); err != nil {
		log.Fatal(err)
	}
	if usbClassStmt, err = db.Prepare(usbClassSQL); err != nil {
		log.Fatal(err)
	}
	if usbSubclassStmt, err = db.Prepare(usbSubclassSQL); err != nil {
		log.Fatal(err)
	}
	if usbProtocolStmt, err = db.Prepare(usbProtocolSQL); err != nil {
		log.Fatal(err)
	}
}

func main() {

	resp, err := http.Get(`http://www.linux-usb.org/usb.ids`)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	var vid, cid, sid string

	for scanner.Scan() {

		if m := vidRgx.FindStringSubmatch(scanner.Text()); m != nil {
			vid = m[1]
			if _, err := usbVendorStmt.Exec(m[1], m[2]); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if m := pidRgx.FindStringSubmatch(scanner.Text()); m != nil {
			fmt.Println(m[1], m[2])
			if _, err := usbProductStmt.Exec(vid, m[1], m[2]); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if m := classRgx.FindStringSubmatch(scanner.Text()); m != nil {
			cid = m[1]
			if _, err := usbClassStmt.Exec(m[1], m[2]); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if m := subclassRgx.FindStringSubmatch(scanner.Text()); m != nil {
			sid = m[1]
			if _, err := usbSubclassStmt.Exec(cid, m[1], m[2]); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if m := protocolRgx.FindStringSubmatch(scanner.Text()); m != nil {
			if _, err := usbSubclassStmt.Exec(cid, sid, m[1], m[2]); err != nil {
				log.Fatal(err)
			}
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

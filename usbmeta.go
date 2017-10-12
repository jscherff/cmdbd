package main

import (
	`bufio`
	`encoding/json`
	`fmt`
	`io/ioutil`
	`net/http`
	`regexp`
	`strconv`

	`github.com/google/gousb`
)

const (
	UsbMetaFileMode = 0640
)

var (
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

// UsbMeta contains all known information USB this.Vendors, products, this.Classes,
// subthis.Classes, and protocols.
type UsbMeta struct {
	Vendors map[gousb.ID]*Vendor
	Classes map[gousb.Class]*Class
}

// Vendor contains the vendor name and mappings to all the vendor's products.
type Vendor struct {
	Name string
	Product	map[gousb.ID]*Product
}

// String returns the name of the vendor.
func (this *Vendor) String() string {
	return this.Name
}

// Product contains the name of the product.
type Product struct {
	Name string
}

// String returns the name of the product.
func (this *Product) String() string {
	return this.Name
}

// Class contains the name of the class and mappings for each subclass.
type Class struct {
	Name string
	Subclass map[gousb.Class]*Subclass
}

// String returns the name of the class.
func (this *Class) String() string {
	return this.Name
}

// Subclass contains the name of the subclass and any associated protocols.
type Subclass struct {
	Name string
	Protocol map[gousb.Protocol]*Protocol
}

// String returns the name of the Subclass.
func (this *Subclass) String() string {
	return this.Name
}

// Protocol contains the name of the protocol.
type Protocol struct {
	Name string
}

// String returns the name of the protocol.
func (this *Protocol) String() string {
	return this.Name
}

// NewUsbMeta creates a new instance of UsbMeta with empty vendor/class maps.
func (this *UsbMeta) Init(refresh bool) (err error) {

	this = &UsbMeta{
		Vendors: make(map[gousb.ID]*Vendor),
		Classes: make(map[gousb.Class]*Class),
	}

	if refresh {
		if err = this.LoadUrl(conf.URLs.UsbMeta); err != nil {
			return err
		}
		if err = this.Save(conf.Files.UsbMeta); err != nil {
			return err
		}
	} else {
		err = this.Load(conf.Files.UsbMeta)
	}

	return err
}

// GetVendor returns the USB vendor associated with a vendor ID.
func (this *UsbMeta) GetVendor(svid string) (*Vendor, error) {

	if vid, err := strconv.ParseUint(svid, 16, 16); err != nil {
		return nil, err
	} else if v, ok := this.Vendors[gousb.ID(vid)]; !ok {
		return nil, fmt.Errorf(`vendor %q not found`, svid)
	} else {
		return v, nil
	}
}

// GetProduct returns the USB product associated with a product ID.
func (this *Vendor) GetProduct(spid string) (*Product, error) {

	if pid, err := strconv.ParseUint(spid, 16, 16); err != nil {
		return nil, err
	} else if p, ok := this.Product[gousb.ID(pid)]; !ok {
		return nil, fmt.Errorf(`product %q not found`, spid)
	} else {
		return p, nil
	}
}

// GetClass returns the USB class associated with a class ID.
func (this *UsbMeta) GetClass(scid string) (*Class, error) {

	if cid, err := strconv.ParseUint(scid, 16, 16); err != nil {
		return nil, err
	} else if c, ok := this.Classes[gousb.Class(cid)]; !ok {
		return nil, fmt.Errorf(`class %q not found`, scid)
	} else {
		return c, nil
	}
}

// GetSubclass returns the USB subclass associated with a subclass ID.
func (this *Class) GetSubclass(ssid string) (*Subclass, error) {

	if sid, err := strconv.ParseUint(ssid, 16, 16); err != nil {
		return nil, err
	} else if s, ok := this.Subclass[gousb.Class(sid)]; !ok {
		return nil, fmt.Errorf(`subclass %q not found`, ssid)
	} else {
		return s, nil
	}
}

// GetProtocol returns the USB protocol associated with a protocol ID.
func (this *Subclass) GetProtocol(spid string) (*Protocol, error) {

	if pid, err := strconv.ParseUint(spid, 16, 16); err != nil {
		return nil, err
	} else if p, ok := this.Protocol[gousb.Protocol(pid)]; !ok {
		return nil, fmt.Errorf(`protocol %q not found`, spid)
	} else {
		return p, nil
	}
}

// LoadUrl reads vendor, product, class, subclass, and protocol information
// line-by-line from an io.Reader source and populates data structures for
// use by the application.
func (this *UsbMeta) LoadUrl(url string) error {

	var (
		vendor   *Vendor
		product  *Product
		class    *Class
		subclass *Subclass
		protocol *Protocol
	)

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {

		if matches := productRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			id, err := strconv.ParseUint(matches[1], 16, 16)

			if err != nil {
				return err
			}
			if vendor == nil {
				return fmt.Errorf(`product with no vendor`)
			}
			if vendor.Product == nil {
				vendor.Product = make(map[gousb.ID]*Product)
			}

			product = &Product{Name: matches[2]}
			vendor.Product[gousb.ID(id)] = product
			continue
		}

		if matches := vendorRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			id, err := strconv.ParseUint(matches[1], 16, 16)

			if err != nil {
				return err
			}

			vendor = &Vendor{Name: matches[2]}
			this.Vendors[gousb.ID(id)] = vendor
			continue
		}

		if matches := protocolRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			id, err := strconv.ParseUint(matches[1], 16, 16)

			if err != nil {
				return err
			}
			if subclass == nil {
				return fmt.Errorf(`protocol with no subclass`)
			}
			if subclass.Protocol == nil {
				subclass.Protocol = make(map[gousb.Protocol]*Protocol)
			}

			protocol = &Protocol{Name: matches[2]}
			subclass.Protocol[gousb.Protocol(id)] = protocol
			continue
		}

		if matches := subclassRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			id, err := strconv.ParseUint(matches[1], 16, 16)

			if err != nil {
				return err
			}
			if class == nil {
				return fmt.Errorf(`subclass with no class`)
			}
			if class.Subclass == nil {
				class.Subclass = make(map[gousb.Class]*Subclass)
			}

			subclass = &Subclass{Name: matches[2]}
			class.Subclass[gousb.Class(id)] = subclass
			continue
		}

		if matches := classRgx.FindStringSubmatch(scanner.Text()); matches != nil {

			id, err := strconv.ParseUint(matches[1], 16, 16)

			if err != nil {
				return err
			}

			class = &Class{Name: matches[2]}
			this.Classes[gousb.Class(id)] = class
		}
	}

	return scanner.Err()
}

// Load loads previously-saved USB information from disk.
func (this *UsbMeta) Load(f string) error {

	j, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}
	if err := json.Unmarshal(j, &this); err != nil {
		return err
	}

	return nil
}

// Save saves current USB information to disk.
func (this *UsbMeta) Save(f string) error {

	j, err := json.Marshal(this)

	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(f, j, UsbMetaFileMode); err != nil {
		return err
	}

	return nil
}

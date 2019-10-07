package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"log"
	"net/http"
)

const (
	ProductTypePort = "MEGAPORT"
	ProductTypeMCR1 = "MEGAPORT"
	ProductTypeMCR2 = "MCR2"
)

type PortsService struct {
	c *Client
}

func NewPortsService(c *Client) *PortsService {
	return &PortsService{c}
}

func (p *PortsService) Create(name string, locationId, speed, term uint64, validate bool) (string, error) {
	payload, err := json.Marshal([]OrderPort{
		OrderPort{
			// CreateDate   uint64          // TODO: need to fill in? :o
			//LagPortCount uint64          // TODO: Required: the number of ports in this LAG order (https://dev.megaport.com/#standard-api-orders-validate-lag-order)
			LocationId: locationId,
			// LocationUid  string // TODO: null in example, is it a string? https://dev.megaport.com/#standard-api-orders-validate-port-order
			// Market : l.Market,
			PortSpeed:   speed,
			ProductName: name,
			ProductType: ProductTypePort,
			Term:        term,
			Virtual:     false,
		},
	})
	if err != nil {
		return "", err
	}
	b := bytes.NewReader(payload)
	if validate { // TODO: do we really want to make this conditional?
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/validate", p.c.BaseURL), b)
		if err != nil {
			return "", err
		}
		if err := p.c.do(req, nil); err != nil {
			return "", err
		}
		b.Seek(0, 0) // TODO: ?
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/buy", p.c.BaseURL), b)
	if err != nil {
		return "", err
	}
	d := []map[string]interface{}{}
	if err := p.c.do(req, &d); err != nil {
		return "", err
	}
	return d[0]["technicalServiceUid"].(string), nil
}

func (p *PortsService) Get(uid string) (*Product, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/product/%s", p.c.BaseURL, uid), nil)
	if err != nil {
		return nil, err
	}
	data := &Product{}
	if err := p.c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (p *PortsService) Update() {}
func (p *PortsService) Delete() {}

// port: virtual = false, type = MEGAPORT
// mcr1: virtual = true , type = MEGAPORT
// mcr2: virtual = false, type = MCR2
type OrderPort struct {
	CreateDate   uint64 `json:"createDate,omitempty"`   // TODO: need to fill in? :o
	LagPortCount uint64 `json:"lagPortCount,omitempty"` // TODO: Required: the number of ports in this LAG order (https://dev.megaport.com/#standard-api-orders-validate-lag-order)
	LocationId   uint64 `json:"locationId"`
	LocationUid  string `json:"locationUid,omitempty"` // TODO: null in example, is it a string? https://dev.megaport.com/#standard-api-orders-validate-port-order
	Market       string `json:"market,omitempty"`      // TODO: what is this ???
	PortSpeed    uint64 `json:"portSpeed"`             // TODO: validate 1000, 10000, 100000
	ProductName  string `json:"productName"`
	ProductType  string `json:"productType"` // TODO: "MEGAPORT"?
	Term         uint64 `json:"term"`
	Virtual      bool   `json:"virtual"` // TODO: False for port, true for MCR1.0 (https://dev.megaport.com/#standard-api-orders-validate-port-order)
}

type OrderPortConfig struct {
	McrASN uint64 `json:"mcrAsn,omitempty"`
}

type OrderVxc struct {
	ProductUid     string                   `json:"productUid"`
	AssociatedVxcs []OrderVxcAssociatedVxcs `json:"associatedVxcs"`
}

type OrderVxcAssociatedVxcs struct {
	ProductName string       `json:"productName"`
	RateLimit   uint64       `json:"rateLimit"`
	AEnd        *OrderVxcEnd `json:"aEnd,omitempty"`
	BEnd        *OrderVxcEnd `json:"bEnd"`
}

type OrderVxcEnd struct {
	ProductUid string `json:"productUid"`
	VLan       uint64 `json:"vlan,omitempty"`
}

func (c *Client) PostPortOrder(o OrderPort) (string, error) {
	p, err := json.Marshal([]OrderPort{o})
	if err != nil {
		return "", err
	}
	b := bytes.NewBuffer(p)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/buy", c.BaseURL), b)
	if err != nil {
		return "", err
	}
	d := []map[string]interface{}{}
	if err := c.do(req, &d); err != nil {
		return "", err
	}
	return d[0]["technicalServiceUid"].(string), nil
}

func (c *Client) PostVxcOrder(o []OrderVxc) error {
	if o == nil {
		o = []OrderVxc{
			OrderVxc{
				ProductUid: "94ec5655-5bef-4734-b172-97f3aed05382",
				AssociatedVxcs: []OrderVxcAssociatedVxcs{
					OrderVxcAssociatedVxcs{
						ProductName: "bar",
						RateLimit:   100,
						BEnd: &OrderVxcEnd{
							ProductUid: "f2c5b25b-e202-4708-9c25-1130c94689b3",
							VLan:       99,
						},
					},
				},
			},
		}
	}
	s, err := json.MarshalIndent(o, "", "  ")
	fmt.Printf("%+v\n", err)
	fmt.Printf("%s\n", s)
	b := bytes.NewBuffer(s)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/buy", c.BaseURL), b)
	if err != nil {
		return nil
	}
	return c.do(req, nil)
}

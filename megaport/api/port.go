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

// port: virtual = false, type = MEGAPORT
// mcr1: virtual = true , type = MEGAPORT
// mcr2: virtual = false, type = MCR2

type PortService struct {
	c *Client
}

func NewPortService(c *Client) *PortService {
	return &PortService{c}
}

func (p *PortService) Create(name string, locationId, speed, term uint64, validate bool) (string, error) {
	payload, err := json.Marshal([]portOrder{portOrder{
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
	}})
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

func (p *PortService) Get(uid string) (*Product, error) {
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

func (p *PortService) Update() error {
	return nil
}

func (p *PortService) Delete(uid string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/product/%s/action/CANCEL_NOW", p.c.BaseURL, uid), nil)
	if err != nil {
		return err
	}
	if err := p.c.do(req, nil); err != nil {
		return err
	}
	return nil
}

func (p *PortService) List() ([]*Product, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/products", p.c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	data := []*Product{}
	if err := p.c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

type portOrder struct {
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

type portOrderConfig struct {
	McrASN uint64 `json:"mcrAsn,omitempty"`
}

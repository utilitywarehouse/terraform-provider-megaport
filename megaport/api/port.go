package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ProductTypePort = "MEGAPORT"
	ProductTypeMCR1 = "MEGAPORT"
	ProductTypeMCR2 = "MCR2"
	ProductTypeVXC  = "VXC"
)

// port: virtual = false, type = MEGAPORT
// mcr1: virtual = true , type = MEGAPORT
// mcr2: virtual = false, type = MCR2

type portCreatePayload struct {
	CreateDate   *uint64 `json:"createDate,omitempty"`   // TODO: need to fill in? :o
	LagPortCount *uint64 `json:"lagPortCount,omitempty"` // TODO: Required: the number of ports in this LAG order (https://dev.megaport.com/#standard-api-orders-validate-lag-order)
	LocationId   *uint64 `json:"locationId"`
	LocationUid  *string `json:"locationUid,omitempty"` // TODO: null in example, is it a string? https://dev.megaport.com/#standard-api-orders-validate-port-order
	Market       *string `json:"market,omitempty"`      // TODO: what is this ???
	PortSpeed    *uint64 `json:"portSpeed"`             // TODO: validate 1000, 10000, 100000
	ProductName  *string `json:"productName"`
	ProductType  *string `json:"productType"` // TODO: "MEGAPORT"?
	Term         *uint64 `json:"term"`
	Virtual      *bool   `json:"virtual"` // TODO: False for port, true for MCR1.0 (https://dev.megaport.com/#standard-api-orders-validate-port-order)
}

type PortCreateInput struct {
	LocationId *uint64
	Name       *string
	Speed      *uint64
	Term       *uint64
}

func (v *PortCreateInput) productType() string {
	return ProductTypePort
}

func (v *PortCreateInput) toPayload() ([]byte, error) {
	payload := []*portCreatePayload{{
		LocationId:  v.LocationId,
		PortSpeed:   v.Speed,
		ProductName: v.Name,
		ProductType: String(ProductTypePort), // TODO
		Term:        v.Term,
		Virtual:     Bool(false), // TODO
	}}
	return json.Marshal(payload)
}

func (c *Client) CreatePort(v *PortCreateInput) (*string, error) {
	d, err := c.create(v)
	if err != nil {
		return nil, err
	}
	uid := d[0]["technicalServiceUid"].(string)
	return &uid, nil
}

func (c *Client) GetPort(uid string) (*Product, error) {
	d := &Product{}
	err := c.get(uid, d)
	return d, err
}

func (c *Client) UpdatePort() error {
	return nil
}

func (c *Client) DeletePort(uid string) error {
	return c.delete(uid)
}

func (c *Client) ListPorts() ([]*Product, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/products", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	data := []*Product{}
	if err := c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

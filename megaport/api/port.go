package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ProductTypePort = "MEGAPORT" // Virtual = false, ProductType = MEGAPORT
	ProductTypeMcr1 = "MEGAPORT" // Virtual = true,  ProductType = MEGAPORT
	ProductTypeMcr2 = "MCR2"     // Virtual = false, ProductType = MCR2
	ProductTypeVxc  = "VXC"
)

// port: virtual = false, type = MEGAPORT
// mcr1: virtual = true , type = MEGAPORT
// mcr2: virtual = false, type = MCR2

type portCreatePayload struct {
	CreateDate            *uint64                      `json:"createDate,omitempty"` // TODO: need to fill in? :o
	Config                *portCreatePayloadPortConfig `json:"config,omitempty"`
	CostCentre            *string                      `json:"costCentre"`
	LagPortCount          *uint64                      `json:"lagPortCount,omitempty"` // TODO: Required: the number of ports in this LAG order (https://dev.megaport.com/#standard-api-orders-validate-lag-order)
	LocationId            *uint64                      `json:"locationId"`
	LocationUid           *string                      `json:"locationUid,omitempty"` // TODO: null in example, is it a string? https://dev.megaport.com/#standard-api-orders-validate-port-order
	Market                *string                      `json:"market,omitempty"`      // TODO: what is this ???
	PortSpeed             *uint64                      `json:"portSpeed"`             // TODO: validate 1000, 10000, 100000
	ProductName           *string                      `json:"productName"`
	ProductType           *string                      `json:"productType"` // TODO: "MEGAPORT"?
	Term                  *uint64                      `json:"term"`
	Virtual               *bool                        `json:"virtual"` // TODO: False for port, true for MCR1.0 (https://dev.megaport.com/#standard-api-orders-validate-port-order)
	MarketplaceVisibility *bool                        `json:"marketplaceVisibility,omitempty"`
}

type portCreatePayloadPortConfig struct {
	McrAsn *uint64 `json:"mcrAsn,omitempty"`
}

type portUpdatePayload struct {
	Name                  *string `json:"name,omitempty"`
	CostCentre            *string `json:"costCentre,omitempty"`
	MarketplaceVisibility *bool   `json:"marketplaceVisibility,omitempty"`
	// RateLimit             *uint64 `json:"rateLimit,omitempty"` // Only applicable to MCR. Must be one of 100, 500, 1000, 2000, 3000, 4000, 5000
}

type PortCreateInput struct {
	LocationId            *uint64
	MarketplaceVisibility *bool
	Name                  *string
	Speed                 *uint64
	Term                  *uint64
	InvoiceReference      *string
}

func (v *PortCreateInput) productType() string {
	return ProductTypePort
}

func (v *PortCreateInput) toPayload() ([]byte, error) {
	payload := []*portCreatePayload{{
		LocationId:            v.LocationId,
		CostCentre:            v.InvoiceReference,
		PortSpeed:             v.Speed,
		ProductName:           v.Name,
		ProductType:           String(ProductTypePort), // TODO
		Term:                  v.Term,
		Virtual:               Bool(false), // TODO
		MarketplaceVisibility: v.MarketplaceVisibility,
	}}
	return json.Marshal(payload)
}

type PortUpdateInput struct {
	InvoiceReference      *string
	MarketplaceVisibility *bool
	Name                  *string
	ProductUid            *string
}

func (v *PortUpdateInput) productType() string {
	return ProductTypePort
}

func (v *PortUpdateInput) toPayload() ([]byte, error) {
	payload := &portUpdatePayload{
		Name:                  v.Name,
		CostCentre:            v.InvoiceReference,
		MarketplaceVisibility: v.MarketplaceVisibility,
	}
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
	if err := c.get(uid, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *Client) UpdatePort(v *PortUpdateInput) error {
	return c.update(*v.ProductUid, v)
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

func (c *Client) GetPortVlanIdAvailable(uid string, vlanId uint64) (bool, error) {
	v := url.Values{}
	v.Set("vlan", strconv.FormatUint(vlanId, 10))
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/product/port/%s/vlan?%s", c.BaseURL, uid, v.Encode()), nil)
	if err != nil {
		return false, err
	}
	data := []uint64{}
	if err := c.do(req, &data); err != nil {
		return false, err
	}
	for _, id := range data {
		if vlanId == id {
			return true, nil
		}
	}
	return false, nil
}

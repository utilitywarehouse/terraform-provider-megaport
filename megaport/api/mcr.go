package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type McrCreateInput interface {
	networkDesignInput
	mcrVersion() uint64
}

type Mcr1CreateInput struct {
	Asn              *uint64
	InvoiceReference *string
	LocationId       *uint64
	Name             *string
	RateLimit        *uint64
	Term             *uint64
}

func (v *Mcr1CreateInput) mcrVersion() uint64 {
	return 1
}

func (v *Mcr1CreateInput) productType() string {
	return ProductTypeMcr1
}

func (v *Mcr1CreateInput) toPayload() ([]byte, error) {
	payload := []*portCreatePayload{{
		LocationId:  v.LocationId,
		Config:      &portCreatePayloadPortConfig{},
		CostCentre:  v.InvoiceReference,
		PortSpeed:   v.RateLimit,
		ProductName: v.Name,
		ProductType: String(ProductTypeMcr1),
		Term:        v.Term,
		Virtual:     Bool(true),
	}}
	if v.Asn != nil && *v.Asn > 0 {
		payload[0].Config.McrAsn = v.Asn
	}
	return json.Marshal(payload)
}

type Mcr2CreateInput struct {
	Asn              *uint64
	InvoiceReference *string
	LocationId       *uint64
	Name             *string
	RateLimit        *uint64
}

func (v *Mcr2CreateInput) mcrVersion() uint64 {
	return 2
}

func (v *Mcr2CreateInput) productType() string {
	return ProductTypeMcr2
}

func (v *Mcr2CreateInput) toPayload() ([]byte, error) {
	payload := []*portCreatePayload{{
		LocationId:  v.LocationId,
		Config:      &portCreatePayloadPortConfig{},
		CostCentre:  v.InvoiceReference,
		PortSpeed:   v.RateLimit,
		ProductName: v.Name,
		ProductType: String(ProductTypeMcr2),
	}}
	if v.Asn != nil && *v.Asn > 0 {
		payload[0].Config.McrAsn = v.Asn
	}
	return json.Marshal(payload)
}

type McrUpdateInput interface {
	networkDesignInput
	productUid() string
}

type Mcr1UpdateInput struct {
	InvoiceReference *string
	Name             *string
	ProductUid       *string
}

func (v *Mcr1UpdateInput) productUid() string {
	return *v.ProductUid
}

func (v *Mcr1UpdateInput) productType() string {
	return ProductTypeMcr1
}

func (v *Mcr1UpdateInput) toPayload() ([]byte, error) {
	payload := &portUpdatePayload{
		Name:       v.Name,
		CostCentre: v.InvoiceReference,
	}
	return json.Marshal(payload)
}

type Mcr2UpdateInput struct {
	InvoiceReference *string
	Name             *string
	ProductUid       *string
}

func (v *Mcr2UpdateInput) productUid() string {
	return *v.ProductUid
}

func (v *Mcr2UpdateInput) productType() string {
	return ProductTypeMcr2
}

func (v *Mcr2UpdateInput) toPayload() ([]byte, error) {
	payload := &portUpdatePayload{
		Name:       v.Name,
		CostCentre: v.InvoiceReference,
	}
	return json.Marshal(payload)
}

func (c *Client) CreateMcr(v McrCreateInput) (*string, error) {
	d, err := c.create(v)
	if err != nil {
		return nil, err
	}
	uid := d[0]["technicalServiceUid"].(string)
	return &uid, nil
}

func (c *Client) GetMcr(uid string) (*Product, error) {
	d := &Product{}
	if err := c.get(uid, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *Client) UpdateMcr(v McrUpdateInput) error {
	return c.update(v.productUid(), v)
}

func (c *Client) DeleteMcr(uid string) error {
	return c.delete(uid)
}

func (c *Client) ListMcrs() ([]*Product, error) {
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

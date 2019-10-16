package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"log"
	"net/http"
)

type VxcService struct {
	c *Client
}

func NewVxcService(c *Client) *VxcService {
	return &VxcService{c}
}

type networkDesignInput interface {
	toPayload() ([]byte, error)
}

func (p *VxcService) create(v networkDesignInput) (*string, error) {
	payload, err := v.toPayload()
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(payload)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/validate", p.c.BaseURL), b)
	if err != nil {
		return nil, err
	}
	if err := p.c.do(req, nil); err != nil {
		return nil, err
	}
	b.Seek(0, 0) // TODO: error handling ?
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/buy", p.c.BaseURL), b)
	if err != nil {
		return nil, err
	}
	d := []map[string]interface{}{}
	if err := p.c.do(req, &d); err != nil {
		return nil, err
	}
	uid := d[0]["vxcJTechnicalServiceUid"].(string)
	return &uid, nil
}

func (p *VxcService) get(uid string) (*ProductAssociatedVxc, error) { // TODO: change name
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/product/%s", p.c.BaseURL, uid), nil)
	if err != nil {
		return nil, err
	}
	data := &ProductAssociatedVxc{}
	if err := p.c.do(req, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (p *VxcService) update(uid string, v networkDesignInput) error {
	payload, err := v.toPayload()
	if err != nil {
		return err
	}
	b := bytes.NewReader(payload)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v2/product/vxc/%s", p.c.BaseURL, uid), b)
	if err != nil {
		return err
	}
	if err := p.c.do(req, nil); err != nil {
		return err
	}
	return nil
}

func (p *VxcService) delete(uid string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/product/%s/action/CANCEL_NOW", p.c.BaseURL, uid), nil)
	if err != nil {
		return err
	}
	if err := p.c.do(req, nil); err != nil {
		return err
	}
	return nil
}

type vxcCreatePayload struct {
	ProductUid     string                           `json:"productUid"`
	AssociatedVxcs []*vxcCreatePayloadAssociatedVxc `json:"associatedVxcs"`
}

type vxcCreatePayloadAssociatedVxc struct {
	ProductName string      `json:"productName"`
	RateLimit   uint64      `json:"rateLimit"`
	CostCentre  string      `json:"costCentre"`
	AEnd        interface{} `json:"aEnd"`
	BEnd        interface{} `json:"bEnd"`
}

type vxcCreatePayloadPrivateVxcEnd struct {
	Vlan uint64 `json:"vlan"`
}

type vxcCreatePayloadPartnerVxcEnd struct {
	ProductUid string `json:"productUid"`
}
type PrivateVxcCreateInput struct {
	InvoiceReference string
	Name             string
	ProductUidA      string
	ProductUidB      string
	RateLimit        uint64
	VlanA            uint64
	VlanB            uint64
}

func (v *PrivateVxcCreateInput) toPayload() ([]byte, error) {
	payload := []*vxcCreatePayload{{
		ProductUid: v.ProductUidA,
		AssociatedVxcs: []*vxcCreatePayloadAssociatedVxc{{
			ProductName: v.Name,
			RateLimit:   v.RateLimit,
			CostCentre:  v.InvoiceReference,
			AEnd:        &vxcCreatePayloadPrivateVxcEnd{Vlan: v.VlanA},
			BEnd:        &vxcCreatePayloadPrivateVxcEnd{ProductUid: v.ProductUidB, Vlan: v.VlanB},
		}},
	}}
	return json.Marshal(payload)
}

type PrivateVxcUpdateInput struct {
	InvoiceReference string
	Name             string
	ProductUid       string
	RateLimit        uint64
	VlanA            uint64
	VlanB            uint64
}

func (v *PrivateVxcUpdateInput) toPayload() ([]byte, error) {
	payload := &vxcUpdatePayload{
		AEndVlan:   v.VlanA,
		BEndVlan:   v.VlanB,
		CostCentre: v.InvoiceReference,
		Name:       v.Name,
		RateLimit:  v.RateLimit,
	}
	return json.Marshal(payload)
}

type vxcUpdatePayload struct {
	AEndVlan   uint64 `json:"aEndVlan"`
	BEndVlan   uint64 `json:"bEndVlan,omitempty"`
	CostCentre string `json:"costCentre"`
	Name       string `json:"name"`
	RateLimit  uint64 `json:"rateLimit"`
}

func (p *VxcService) CreatePrivateVxc(v *PrivateVxcCreateInput) (*string, error) {
	return p.create(v)
}

func (p *VxcService) GetPrivateVxc(uid string) (*ProductAssociatedVxc, error) {
	return p.get(uid)
}

func (p *VxcService) UpdatePrivateVxc(v *PrivateVxcUpdateInput) error {
	return p.update(v.ProductUid, v)
}

func (p *VxcService) DeletePrivateVxc(uid string) error {
	return p.delete(uid)
}

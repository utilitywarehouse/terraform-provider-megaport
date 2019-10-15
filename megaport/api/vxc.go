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

type VxcCreateInput struct {
	InvoiceReference string
	Name             string
	ProductUidA      string
	ProductUidB      string
	RateLimit        uint64
	VlanA            uint64
	VlanB            uint64
}

func (v *VxcCreateInput) toPayload() ([]byte, error) {
	payload := []*vxcCreatePayload{{
		ProductUid: v.ProductUidA,
		AssociatedVxcs: []vxcCreatePayloadAssociatedVxc{{
			ProductName: v.Name,
			RateLimit:   v.RateLimit,
			CostCentre:  v.InvoiceReference,
			AEnd:        &vxcCreatePayloadAssociatedVxcEnd{Vlan: v.VlanA},
			BEnd:        &vxcCreatePayloadAssociatedVxcEnd{ProductUid: v.ProductUidB, Vlan: v.VlanB},
		}},
	}}
	return json.Marshal(payload)
}

type VxcCreateOutput struct {
	ProductUid string
}

type vxcCreatePayload struct {
	ProductUid     string                          `json:"productUid"`
	AssociatedVxcs []vxcCreatePayloadAssociatedVxc `json:"associatedVxcs"`
}

type vxcCreatePayloadAssociatedVxc struct {
	ProductName string                            `json:"productName"`
	RateLimit   uint64                            `json:"rateLimit"`
	CostCentre  string                            `json:"costCentre"`
	AEnd        *vxcCreatePayloadAssociatedVxcEnd `json:"aEnd"`
	BEnd        *vxcCreatePayloadAssociatedVxcEnd `json:"bEnd"`
}

type vxcCreatePayloadAssociatedVxcEnd struct {
	ProductUid string `json:"productUid,omitempty"`
	Vlan       uint64 `json:"vlan"`
}

type VxcUpdateInput struct {
	InvoiceReference string
	Name             string
	ProductUid       string
	RateLimit        uint64
	VlanA            uint64
	VlanB            uint64
}

func (v *VxcUpdateInput) toPayload() ([]byte, error) {
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

func (p *VxcService) Create(v *VxcCreateInput) (*VxcCreateOutput, error) {
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
	return &VxcCreateOutput{ProductUid: d[0]["vxcJTechnicalServiceUid"].(string)}, nil
}

func (p *VxcService) Get(uid string) (*ProductAssociatedVxc, error) {
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

func (p *VxcService) Update(v *VxcUpdateInput) error {
	payload, err := v.toPayload()
	if err != nil {
		return err
	}
	b := bytes.NewReader(payload)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v2/product/vxc/%s", p.c.BaseURL, v.ProductUid), b)
	if err != nil {
		return err
	}
	if err := p.c.do(req, nil); err != nil {
		return err
	}
	return nil
}

func (p *VxcService) Delete(uid string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/product/%s/action/CANCEL_NOW", p.c.BaseURL, uid), nil)
	if err != nil {
		return err
	}
	if err := p.c.do(req, nil); err != nil {
		return err
	}
	return nil
}

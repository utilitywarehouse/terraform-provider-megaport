package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"log"
	"net/http"
)

type vxcOrder struct {
	ProductUid     string                   `json:"productUid"`
	AssociatedVxcs []vxcOrderAssociatedVxcs `json:"associatedVxcs"`
}

type vxcOrderAssociatedVxcs struct {
	ProductName string       `json:"productName"`
	RateLimit   uint64       `json:"rateLimit"`
	AEnd        *vxcOrderEnd `json:"aEnd,omitempty"`
	BEnd        *vxcOrderEnd `json:"bEnd"`
}

type vxcOrderEnd struct {
	ProductUid string `json:"productUid"`
	VLan       uint64 `json:"vlan,omitempty"`
}

type VxcService struct {
	c *Client
}

func NewVxcService(c *Client) *VxcService {
	return &VxcService{c}
}

func (p *VxcService) Create(productAUid, productBUid, name string, vlanA, vlanB, rateLimit uint64) (string, error) {
	order := []vxcOrder{vxcOrder{
		ProductUid: productAUid,
		AssociatedVxcs: []vxcOrderAssociatedVxcs{
			vxcOrderAssociatedVxcs{
				ProductName: name,
				RateLimit:   rateLimit,
				BEnd: &vxcOrderEnd{
					ProductUid: productBUid,
				},
			},
		},
	}}
	if vlanA != 0 {
		order[0].AssociatedVxcs[0].AEnd = &vxcOrderEnd{VLan: vlanB}
	}
	if vlanB != 0 {
		order[0].AssociatedVxcs[0].BEnd.VLan = vlanB
	}
	payload, err := json.Marshal(order)
	if err != nil {
		return "", err
	}
	b := bytes.NewReader(payload)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/validate", p.c.BaseURL), b)
	if err != nil {
		return "", err
	}
	if err := p.c.do(req, nil); err != nil {
		return "", err
	}
	b.Seek(0, 0) // TODO: error handling ?
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/buy", p.c.BaseURL), b)
	if err != nil {
		return "", err
	}
	d := []map[string]interface{}{}
	if err := p.c.do(req, &d); err != nil {
		return "", err
	}
	return d[0]["vxcJTechnicalServiceUid"].(string), nil
}

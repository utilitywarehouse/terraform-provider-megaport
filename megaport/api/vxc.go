package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	//	"log"
	"net/http"
)

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

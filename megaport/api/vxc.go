package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	vxcConnectTypeAws    = "AWS"
	vxcConnectTypeGoogle = "GOOGLE"
)

type networkDesignInput interface {
	toPayload() ([]byte, error)
	productType() string
}

func (c *Client) create(v networkDesignInput) ([]map[string]interface{}, error) {
	payload, err := v.toPayload()
	if err != nil {
		return nil, err
	}
	b := bytes.NewReader(payload)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/validate", c.BaseURL), b)
	if err != nil {
		return nil, err
	}
	if err := c.do(req, nil); err != nil {
		return nil, err
	}
	if _, err := b.Seek(0, 0); err != nil {
		return nil, err
	}
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/networkdesign/buy", c.BaseURL), b)
	if err != nil {
		return nil, err
	}
	d := []map[string]interface{}{}
	if err := c.do(req, &d); err != nil {
		return nil, err
	}
	return d, nil
}

func (c *Client) get(uid string, v interface{}) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v2/product/%s", c.BaseURL, uid), nil)
	if err != nil {
		return err
	}
	if err := c.do(req, v); err != nil {
		return err
	}
	return nil
}

func (c *Client) update(uid string, v networkDesignInput) error {
	payload, err := v.toPayload()
	if err != nil {
		return err
	}
	b := bytes.NewReader(payload)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v2/product/%s/%s", c.BaseURL, strings.ToLower(v.productType()), uid), b)
	if err != nil {
		return err
	}
	if err := c.do(req, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) delete(uid string) error {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/product/%s/action/CANCEL_NOW", c.BaseURL, uid), nil)
	if err != nil {
		return err
	}
	if err := c.do(req, nil); err != nil {
		return err
	}
	return nil
}

type vxcCreatePayload struct {
	ProductUid     *string                          `json:"productUid,omitempty"`
	AssociatedVxcs []*vxcCreatePayloadAssociatedVxc `json:"associatedVxcs,omitempty"`
}

type vxcCreatePayloadAssociatedVxc struct {
	ProductName   *string                 `json:"productName,omitempty"`
	RateLimit     *uint64                 `json:"rateLimit,omitempty"`
	CostCentre    *string                 `json:"costCentre,omitempty"`
	AEnd          *vxcCreatePayloadVxcEnd `json:"aEnd,omitempty"`
	BEnd          *vxcCreatePayloadVxcEnd `json:"bEnd,omitempty"`
	PartnerConfig interface{}             `json:"partnerConfigs,omitempty"`
}

type vxcCreatePayloadVxcEnd struct {
	ProductUid *string `json:"productUid,omitempty"`
	Vlan       *uint64 `json:"vlan,omitempty"`
}

type PrivateVxcCreateInput struct {
	InvoiceReference *string
	Name             *string
	ProductUidA      *string
	ProductUidB      *string
	RateLimit        *uint64
	VlanA            *uint64
	VlanB            *uint64
}

func (v *PrivateVxcCreateInput) productType() string {
	return ProductTypeVxc
}

func (v *PrivateVxcCreateInput) toPayload() ([]byte, error) {
	payload := []*vxcCreatePayload{{ProductUid: v.ProductUidA}}
	av := &vxcCreatePayloadAssociatedVxc{
		ProductName: v.Name,
		RateLimit:   v.RateLimit,
		CostCentre:  v.InvoiceReference,
	}
	if v.VlanA != nil {
		av.AEnd = &vxcCreatePayloadVxcEnd{Vlan: v.VlanA}
	}
	bEnd := &vxcCreatePayloadVxcEnd{ProductUid: v.ProductUidB, Vlan: v.VlanB}
	if *bEnd != (vxcCreatePayloadVxcEnd{}) {
		av.BEnd = bEnd
	}
	if *av != (vxcCreatePayloadAssociatedVxc{}) {
		payload[0].AssociatedVxcs = []*vxcCreatePayloadAssociatedVxc{av}
	}
	return json.Marshal(payload)
}

type vxcUpdatePayload struct {
	AEndVlan   *uint64     `json:"aEndVlan,omitempty"`
	BEndVlan   *uint64     `json:"bEndVlan,omitempty"`
	CostCentre *string     `json:"costCentre,omitempty"`
	Name       *string     `json:"name,omitempty"`
	BEndConfig interface{} `json:"bEndConfig,omitempty"`
	RateLimit  *uint64     `json:"rateLimit,omitempty"`
	// SecondaryName *string     `json:"secondaryName,omitempty"`
}

type PrivateVxcUpdateInput struct {
	InvoiceReference *string
	Name             *string
	ProductUid       *string
	RateLimit        *uint64
	VlanA            *uint64
	VlanB            *uint64
}

func (v *PrivateVxcUpdateInput) productType() string {
	return ProductTypeVxc
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

func (c *Client) CreatePrivateVxc(v *PrivateVxcCreateInput) (*string, error) {
	d, err := c.create(v)
	if err != nil {
		return nil, err
	}
	uid := d[0]["vxcJTechnicalServiceUid"].(string)
	return &uid, nil
}

func (c *Client) GetVxc(uid string) (*ProductAssociatedVxc, error) { // TODO: rename struct
	d := &ProductAssociatedVxc{}
	err := c.get(uid, d)
	return d, err
}

func (c *Client) DeleteVxc(uid string) error {
	return c.delete(uid)
}

func (c *Client) UpdatePrivateVxc(v *PrivateVxcUpdateInput) error {
	return c.update(*v.ProductUid, v)
}

type PartnerConfig interface {
	connectType() string
	toPayload() interface{}
}

type PartnerConfigAws struct {
	AmazonIPAddress   *string
	AwsConnectionName *string
	AwsAccountId      *string
	BGPAuthKey        *string
	CustomerASN       *uint64
	CustomerIPAddress *string
	Type              *string
}

func (v *PartnerConfigAws) connectType() string {
	return "AWS"
}

func (v *PartnerConfigAws) toPayload() interface{} {
	return &vxcCreatePayloadPartnerConfigAws{
		AmazonIpAddress:   v.AmazonIPAddress,
		Asn:               v.CustomerASN,
		AuthKey:           v.BGPAuthKey,
		ConnectType:       String(v.connectType()),
		CustomerIpAddress: v.CustomerIPAddress,
		Name:              v.AwsConnectionName,
		OwnerAccount:      v.AwsAccountId,
		Type:              v.Type,
	}
}

type vxcCreatePayloadPartnerConfigAws struct {
	// AmazonAsn         *uint64 `json:",omitempty"`
	AmazonIpAddress *string `json:"amazonIpAddress,omitempty"`
	Asn             *uint64 `json:"asn,omitempty"`
	AuthKey         *string `json:"authKey,omitempty"`
	ConnectType     *string `json:"connectType,omitempty"`
	// Complete          *bool   `json:"complete,omitempty"`
	CustomerIpAddress *string `json:"customerIpAddress,omitempty"`
	Name              *string `json:"name,omitempty"`
	OwnerAccount      *string `json:"ownerAccount,omitempty"`
	// Prefixes          *string `json:",omitempty"`
	Type *string `json:"type,omitempty"`
}

type CloudVxcCreateInput struct {
	InvoiceReference *string
	Name             *string
	PartnerConfig    PartnerConfig
	ProductUidA      *string
	ProductUidB      *string
	RateLimit        *uint64
	VlanA            *uint64
}

func (v *CloudVxcCreateInput) toPayload() ([]byte, error) {
	payload := []*vxcCreatePayload{{ProductUid: v.ProductUidA}}
	av := &vxcCreatePayloadAssociatedVxc{
		CostCentre:    v.InvoiceReference,
		PartnerConfig: v.PartnerConfig.toPayload(),
		ProductName:   v.Name,
		RateLimit:     v.RateLimit,
	}
	if v.VlanA != nil {
		av.AEnd = &vxcCreatePayloadVxcEnd{Vlan: v.VlanA}
	}
	bEnd := &vxcCreatePayloadVxcEnd{ProductUid: v.ProductUidB}
	if *bEnd != (vxcCreatePayloadVxcEnd{}) {
		av.BEnd = bEnd
	}
	if *av != (vxcCreatePayloadAssociatedVxc{}) {
		payload[0].AssociatedVxcs = []*vxcCreatePayloadAssociatedVxc{av}
	}
	return json.Marshal(payload)
}

func (v *CloudVxcCreateInput) productType() string {
	return ProductTypeVxc
}

type CloudVxcUpdateInput struct {
	InvoiceReference *string
	Name             *string
	ProductUid       *string
	PartnerConfig    PartnerConfig
	RateLimit        *uint64
	VlanA            *uint64
}

func (v *CloudVxcUpdateInput) productType() string {
	return ProductTypeVxc
}

func (v *CloudVxcUpdateInput) toPayload() ([]byte, error) {
	payload := &vxcUpdatePayload{
		AEndVlan:   v.VlanA,
		CostCentre: v.InvoiceReference,
		Name:       v.Name,
		BEndConfig: v.PartnerConfig.toPayload(),
		RateLimit:  v.RateLimit,
	}
	return json.Marshal(payload)
}

func (c *Client) CreateCloudVxc(v *CloudVxcCreateInput) (*string, error) {
	d, err := c.create(v)
	if err != nil {
		return nil, err
	}
	uid := d[0]["vxcJTechnicalServiceUid"].(string)
	return &uid, nil
}

func (c *Client) UpdateCloudVxc(v *CloudVxcUpdateInput) error {
	return c.update(*v.ProductUid, v)
}

package api

import (
	"encoding/json"
)

const (
	VxcTypePrivate = "private"
	VxcTypeAws     = "aws"
	VxcTypePartner = "partner"
)

// Some of the following types differ from examples seen in the documentation at
// https://dev.megaport.com. In some cases, the API responds with a different
// set of fields and/or field types and so the structs in this file match that,
// instead of those in the documentation.

type megaportResponse struct {
	Message string
	Data    interface{}
}

type responseLoginData struct {
	Token string
}

// Location data
type Location struct {
	Address          LocationAddress
	Campus           string
	Country          string
	Id               uint64
	Latitude         float64
	LiveDate         uint64
	Longitude        float64
	Market           string
	Metro            string
	Name             string
	NetworkRegion    string
	Products         LocationProducts
	SiteCode         string
	Status           string
	VRouterAvailable bool
}

// Address data
type LocationAddress struct {
	City     string
	Country  string
	Postcode string
	State    string
	Street   string
	Suburb   string
}

type LocationProducts struct {
	MCR        bool
	MCRVersion uint64
	MCR1       []uint64
	MCR2       []uint64
	Megaport   []uint64
}

type Megaport struct {
	AggregationId uint64 `json:"aggregation_id"`
	CompanyName   string
	CompanyUid    string
	ConnectType   string
	LagId         uint64 `json:"lag_id"`
	LagPrimary    bool   `json:"lag_primary"`
	LocationId    uint64
	ProductUid    string
	Rank          uint64
	Speed         uint64
	Title         string
	VxcPermitted  bool
}

type InternetExchange struct {
	ASN           uint64
	Description   string
	ECIX          bool
	GroupMetro    string `json:"group_metro"`
	Name          string
	NetworkRegion string `json:"network_region"`
	PrimaryIPv4   InternetExchangeIPAddress
	PrimaryIPv6   InternetExchangeIPAddress
	SecondaryIPv4 InternetExchangeIPAddress
	SecondaryIPv6 InternetExchangeIPAddress
	State         string
}

type InternetExchangeIPAddress struct {
	Type  string
	Value string
}

type Product struct {
	AdminLocked bool
	// AggregationId // TODO: haven't seen a value other than null
	// AssociatedIxs []ProductsAssociatedIx // TODO: haven't seen a value other than an empty list
	AssociatedVxcs []ProductAssociatedVxc
	// AttributeTags // TODO: haven't seen a value other than an empty map
	BuyoutPort         bool
	Cancelable         bool
	CompanyName        string
	CompanyUid         string
	ContractStartDate  uint64
	ContractEndDate    uint64
	ContractTermMonths uint64
	CostCentre         string
	CreateDate         uint64
	CreatedBy          string
	// LagId // TODO: haven't seen a value other than null
	LagPrimary            bool
	LiveDate              uint64
	LocationId            uint64
	Locked                bool
	Market                string
	MarketplaceVisibility bool
	PortSpeed             uint64
	ProductName           string
	ProductType           string
	ProductUid            string
	ProvisioningStatus    string
	Resources             ProductResources
	// SecondaryName // TODO: haven't seen a value other than null
	// TerminateDate // TODO: haven't seen a value other than null
	// UsageAlgorithm // TODO: haven't seen a value other than null
	Virtual         bool
	VxcPermitted    bool
	VxcAutoApproval bool
}

type ProductResources struct { // TODO: verify these are the only valid fields
	// CrossConnect  ProductResourcesCrossConnect `json:"cross_connect"` // TODO: only referenced in https://dev.megaport.com/#general-get-product-list
	Interface     ProductResourcesInterface
	VirtualRouter ProductResourcesVirtualRouter `json:"virtual_router"`
	VLL           ProductResourcesVLL
}

type ProductResourcesInterface struct {
	Demarcation  string
	Description  string
	Id           uint64 `json:"-"`
	LoaTemplate  string `json:"loa_template"`
	Media        string
	Name         string
	PortSpeed    uint64 `json:"-"`
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
	// SupportedSpeeds []uint64 `json:"supported_speeds"` // TODO: only referenced in https://dev.megaport.com/#general-get-product-list
	Up uint64 `json:"-"`
}

type productResourcesInterfaceFloats struct {
	Id        float64
	PortSpeed float64 `json:"port_speed"`
	Up        float64
}

type productResourcesInterface ProductResourcesInterface

func (pr *ProductResourcesInterface) UnmarshalJSON(b []byte) (err error) {
	v := productResourcesInterface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*pr = ProductResourcesInterface(v)
	vf := productResourcesInterfaceFloats{}
	if err := json.Unmarshal(b, &vf); err != nil {
		return err
	}
	pr.Id = uint64(vf.Id)
	pr.PortSpeed = uint64(vf.PortSpeed)
	pr.Up = uint64(vf.Up)
	return nil
}

type ProductResourcesVirtualRouter struct {
	Id           uint64 `json:"-"`
	McrASN       uint64 `json:"-"`
	Name         string
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
	Speed        uint64 `json:"-"`
}

type productResourcesVirtualRouterFloats struct {
	Id     float64 `json:"id"`
	McrASN float64 `json:"mcrAsn"`
	Speed  float64 `json:"speed"`
}

type productResourcesVirtualRouter ProductResourcesVirtualRouter

func (pr *ProductResourcesVirtualRouter) UnmarshalJSON(b []byte) (err error) {
	v := productResourcesVirtualRouter{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*pr = ProductResourcesVirtualRouter(v)
	vf := productResourcesVirtualRouterFloats{}
	if err := json.Unmarshal(b, &vf); err != nil {
		return err
	}
	pr.Id = uint64(vf.Id)
	pr.McrASN = uint64(vf.McrASN)
	pr.Speed = uint64(vf.Speed)
	return nil
}

type ProductResourcesVLL struct {
	AVLan        uint64 `json:"-"`
	BVLan        uint64 `json:"-"`
	Description  string
	Id           uint64 `json:"-"`
	Name         string
	RateLimit    uint64 `json:"-"`
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
	Up           uint64 `json:"-"`
}

type productResourcesVLLFloats struct {
	AVLan     float64 `json:"a_vlan"`
	BVLan     float64 `json:"b_vlan"`
	Id        float64 `json:"id"`
	RateLimit float64 `json:"rate_limit_mbps"`
	Up        float64 `json:"up"`
}

type productResourcesVLL ProductResourcesVLL

func (pr *ProductResourcesVLL) UnmarshalJSON(b []byte) (err error) {
	v := productResourcesVLL{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*pr = ProductResourcesVLL(v)
	vf := productResourcesVLLFloats{}
	if err := json.Unmarshal(b, &vf); err != nil {
		return err
	}
	pr.AVLan = uint64(vf.AVLan)
	pr.BVLan = uint64(vf.BVLan)
	pr.Id = uint64(vf.Id)
	pr.RateLimit = uint64(vf.RateLimit)
	pr.Up = uint64(vf.Up)
	return nil
}

type ProductAssociatedVxc struct {
	AdminLocked bool
	// AttributeTags // TODO: haven't seen a value other than an empty map
	AEnd               ProductAssociatedVxcEnd
	BEnd               ProductAssociatedVxcEnd
	Cancelable         bool
	ContractEndDate    uint64 // TODO: haven't seen a value other than null, despite the note in https://dev.megaport.com/#general-get-product-list
	ContractStartDate  uint64 // TODO: haven't seen a value other than null, despite the note in https://dev.megaport.com/#general-get-product-list
	ContractTermMonths uint64
	CostCentre         string
	CreatedBy          string // TODO: haven't seen a value other than null
	CreateDate         uint64
	DistanceBand       string
	Locked             bool
	// NServiceId // TODO: haven't seen a value other than null
	ProductName        string
	ProductType        string
	ProductUid         string
	ProvisioningStatus string
	RateLimit          uint64
	Resources          ProductAssociatedVxcResources // TODO: not documented - is the struct here the same as in Product?
	SecondaryName      string
	UsageAlgorithm     string
	VxcApproval        ProductAssociatedVxcApproval
}

func (v *ProductAssociatedVxc) Type() string {
	if v.AEnd.OwnerUid == v.BEnd.OwnerUid {
		return VxcTypePrivate
	}
	if v.Resources.AwsVirtualInterface.ConnectType == "AWS" {
		return VxcTypeAws
	}
	return VxcTypePartner
}

type ProductAssociatedVxcEnd struct {
	LocationId  uint64
	Location    string
	OwnerUid    string
	ProductUid  string
	ProductName string
	Vlan        uint64
	// SecondaryName // TODO: haven't seen a value other than null
}

type ProductAssociatedVxcApproval struct {
	// Message // TODO: haven't seen a value other than null
	// NewSpeed // TODO: haven't seen a value other than null
	// Status // TODO: haven't seen a value other than null
	// Type // TODO: haven't seen a value other than null
	// Uid // TODO: haven't seen a value other than null
}

type ProductAssociatedVxcResources struct {
	AwsVirtualInterface ProductAssociatedVxcResourcesAwsVirtualInterface `json:"csp_connection"`
}

type ProductAssociatedVxcResourcesAwsVirtualInterface struct {
	Account         string
	AmazonAsn       uint64 `json:"-"`
	AmazonIpAddress string
	AmazonAddress   string `json:"Amazon_address"`
	// Amazon_Asn       uint64 `json:"Amazon_asn"`
	Asn     uint64 `json:"-"`
	AuthKey string
	// Auth_key string `json:"Auth_key"`
	ConnectType       string
	CustomerIpAddress string
	// Customer_address string `json:"Customer_address"`
	Id           uint64 `json:"-"`
	Name         string
	OwnerAccount string
	PeerAsn      uint64 `json:"-"`
	// Prefixes // null?
	ResourceName string `json:"Resource_name"`
	ResourceType string `json:"Resource_type"`
	Type         string
	VifId        string `json:"Vif_id"`
	Vlan         uint64 `json:"-"`
}

type productAssociatedVxcResourcesAwsVirtualInterfaceFloats struct {
	AmazonAsn float64 `json:"amazonAsn"`
	Asn       float64 `json:"asn"`
	Id        float64 `json:"id"`
	PeerAsn   float64 `json:"peerAsn"`
	Vlan      float64 `json:"vlan"`
}

type productAssociatedVxcResourcesAwsVirtualInterface ProductAssociatedVxcResourcesAwsVirtualInterface

func (pr *ProductAssociatedVxcResourcesAwsVirtualInterface) UnmarshalJSON(b []byte) (err error) {
	v := productAssociatedVxcResourcesAwsVirtualInterface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*pr = ProductAssociatedVxcResourcesAwsVirtualInterface(v)
	vf := productAssociatedVxcResourcesAwsVirtualInterfaceFloats{}
	if err := json.Unmarshal(b, &vf); err != nil {
		return err
	}
	pr.AmazonAsn = uint64(vf.AmazonAsn)
	pr.Asn = uint64(vf.Asn)
	pr.Id = uint64(vf.Id)
	pr.PeerAsn = uint64(vf.PeerAsn)
	pr.Vlan = uint64(vf.Vlan)
	return nil
}

type MegaportCharges struct {
	Currency             string
	DailyRate            float64
	DailySetup           float64
	Empty                bool
	FixedRecurringCharge float64
	// ForceProductChange // TODO: haven't seen a value other than null
	HourlyRate  float64
	HourlySetup float64
	// Key string // TODO: haven't seen a value other than "no key"
	LongHaulMbpsRate float64
	MbpsRate         float64
	MonthlyRate      float64
	MonthlySetup     float64
	// PostPaidBaseRate // TODO: haven't seen a value other than "no base rate"
	ProductType string
}

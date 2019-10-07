package api

import (
	"encoding/json"
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
	Address          Address
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
type Address struct {
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
	PrimaryIPv4   IPAddress
	PrimaryIPv6   IPAddress
	SecondaryIPv4 IPAddress
	SecondaryIPv6 IPAddress
	State         string
}

type IPAddress struct {
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
}

type ProductResourcesInterface struct {
	Demarcation  string
	Description  string
	Id           uint64 `json:"-"`
	_id          float64
	LoaTemplate  string `json:"loa_template"`
	Media        string
	Name         string
	PortSpeed    uint64  `json:"-"`
	_portSpeed   float64 `json:"port_speed"`
	ResourceName string  `json:"resource_name"`
	ResourceType string  `json:"resource_type"`
	// SupportedSpeeds []uint64 `json:"supported_speeds"` // TODO: only referenced in https://dev.megaport.com/#general-get-product-list
	Up  uint64 `json:"-"`
	_up float64
}

type productResourcesInterface ProductResourcesInterface

func (pri *ProductResourcesInterface) UnmarshalJSON(b []byte) (err error) {
	v := productResourcesInterface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	*pri = ProductResourcesInterface(v)
	pri.Id = uint64(pri._id)
	pri.PortSpeed = uint64(pri._portSpeed)
	pri.Up = uint64(pri._up)
	return nil
}

type ProductResourcesVirtualRouter struct {
	Id           uint64
	McrASN       uint64
	Name         string
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
	Speed        uint64
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
	SecondaryName      string
	UsageAlgorithm     string
	VxcApproval        ProductAssociatedVxcApproval
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

type Charges struct {
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

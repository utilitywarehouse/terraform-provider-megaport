package api

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

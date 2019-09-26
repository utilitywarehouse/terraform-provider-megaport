package main

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
	Products         Products
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

type Products struct {
	MCR        bool
	MCRVersion uint64
	MCR1       []uint64
	MCR2       []uint64
	Megaport   []uint64
}

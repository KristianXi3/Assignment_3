package entity

type UserRandom struct {
	Id        int     `json:"id"`
	Uid       string  `json:"uid"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Username  string  `json:"username"`
	Address   Address `json:"address"`
}

type Address struct {
	City          string     `json:"city"`
	StreetName    string     `json:"street_name"`
	StreetAddress string     `json:"street_address"`
	ZipCode       string     `json:"zip_code"`
	State         string     `json:"state"`
	Country       string     `json:"country"`
	Coordinates   Coordinate `json:"coordinates"`
}

type Coordinate struct {
	Lattitude float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

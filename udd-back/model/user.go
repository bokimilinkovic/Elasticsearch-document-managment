package model

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type User struct {
	Loc      Location `json:"location"`
	Username string   `json:"username"`
	Country  string   `json:"country"`
	City     string   `json:"city"`
}

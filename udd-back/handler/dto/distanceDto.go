package dto

type SearchByDistanceDto struct {
	City  string `json:"city"`
	Range int    `json:"range"`
}

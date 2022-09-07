package model

type Parameter struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type Weather struct {
	Status Parameter `json:"status"`
}

type WeatherTemplate struct {
	Water  int    `json:"water"`
	Wind   int    `json:"wind"`
	Status string `json:"status"`
}

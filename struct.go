package main

type Geo struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Station struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	IsCoffee int    `json:"is_coffee"`
	IsFood   int    `json:"is_food"`
	Geo      Geo    `json:"geo"`
	Time     int    `json:"time"`
}

type GetStationsListResponse struct {
	Stations []Station `json:"stations"`
}

type GetStationsListRequest struct {
	Train    string `json:"train"`
	Carriage string `json:"carriage,omitempty"`
}

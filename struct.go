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
	Carriage string `json:"carriage"`
}

type OrderItem struct {
	CoffeeType string `json:"coffee_type"`
	Number     int    `json:"number"`
}

type PostOrderRequest struct {
	Train       string      `json:"train"`
	Carriage    string      `json:"carriage"`
	Station     string      `json:"station"`
	RepeatOrder bool        `json:"repeat_order"`
	Delivery    bool        `json:"delivery"`
	Order       []OrderItem `json:"order"`
}

type PostOrderResponse struct {
	ID string `json:"id"`
}

type TaskItem struct {
	Train       string      `json:"train"`
	Carriage    string      `json:"carriage"`
	Station     string      `json:"station"`
	RepeatOrder bool        `json:"repeat_order"`
	Delivery    bool        `json:"delivery"`
	Order       []OrderItem `json:"order"`
	ArrivalTime int64       `json:"arrival_time"`
	CreateTime  int64       `json:"create_time"`
}

type GetTaskListResponse struct {
	Tasks []TaskItem `json:"tasks"`
}

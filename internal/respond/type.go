package respond

import "argc.in/shrt/internal/model"

type msg struct {
	Ok bool `json:"ok"`
}

type msgErr struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type msgRoute struct {
	Ok    bool         `json:"ok"`
	Route *model.Route `json:"route"`
}

type msgRoutes struct {
	Ok     bool          `json:"ok"`
	Routes []model.Route `json:"routes"`
}

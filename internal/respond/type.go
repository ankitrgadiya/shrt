package respond

import "argc.in/shrt/internal/model"

type Msg struct {
	Ok     bool          `json:"ok"`
	Error  string        `json:"error,omitempty"`
	Route  *model.Route  `json:"route,omitempty"`
	Routes []model.Route `json:"routes,omitempty"`
}

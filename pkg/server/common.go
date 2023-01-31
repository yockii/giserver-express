package server

import "github.com/yockii/giserver-express/pkg/database"

type CommonResponse struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type TimeCondition struct {
	Start database.DateTime `json:"start,omitempty" query:"start"`
	End   database.DateTime `json:"end,omitempty" query:"end"`
}

type Paginate struct {
	Total  int64       `json:"total"`
	Offset int         `json:"offset"`
	Limit  int         `json:"limit"`
	Items  interface{} `json:"items"`
}

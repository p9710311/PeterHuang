package models

import "edgea/enums"

// JsonResult 用於返回ajax請求的基類
type JsonResult struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Obj  interface{}          `json:"obj"`
}

// BaseQueryParam 用於查詢的類
type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int    `json:"limit"`
}

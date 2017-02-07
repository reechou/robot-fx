package controller

const (
	ORDER_RESPONSE_OK = iota
	ORDER_ERROR
)

type OrderResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

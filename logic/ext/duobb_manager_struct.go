package ext

import (
	"errors"
)

const (
	DUOBB_MANAGER_RESPONSE_OK              = 1000
	DUOBB_MANAGER_GOODS_SEARCH_NO_DISCOUNT = 2001
)

var (
	ERR_DUOBB_GOODS_SEARCH_NO_DISCOUNT = errors.New("goods has no discount")
)

type GoodsSearchReq struct {
	Query    string `json:"q,omitempty"`
	Alimama  string `json:"userName"`
	Pid      string `json:"pid"`
	AdzoneId string `json:"adzoneid"`
	SiteId   string `json:"siteid"`
	Url      string `json:"url,omitempty"`
	Offset   int64  `json:"offset"`
}

type GoodsSearchData struct {
	Token              string  `json:"token"`
	Title              string  `json:"title"`
	PicUrl             string  `json:"picUrl"`
	EffectiveStartTime string  `json:"effectiveStartTime"`
	EffectiveEndTime   string  `json:"effectiveEndTime"`
	ZkPrice            float32 `json:"zkPrice"`    // 商品原价
	StartFee           float32 `json:"startFee"`   // 优惠满多少
	Amount             float32 `json:"amount"`     // 优惠券金额
	RlRate             float32 `json:"rlRate"`     // 佣金比例
	BestChoose         int     `json:"bestChoose"` // 1表示通用 2表示高佣计划 3表示定向计划
	EndPrice           float32 `json:"endPrice"`   // 券后价
	UlandUrl           string  `json:"ulandUrl"`   // 优惠券地址
	ShortUrl           string  `json:"sUrl"`
}

type GoodsSearchRsp struct {
	Code int64           `json:"state"`
	Msg  string          `json:"message,omitempty"`
	Data GoodsSearchData `json:"data,omitempty"`
}

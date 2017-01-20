package ext

const (
	WECHAT_RESPONSE_OK = 1000
)

const (
	WX_SEND_MSG_CHAN_LEN = 10240
)

type WithdrawalReq struct {
	OpenId      string `json:"openId"`
	TotalAmount int64  `json:"total_amount"`
	MchBillno   string `json:"mch_billno"`
}

type WeixinMsgSendReq struct {
	OpenId    string  `json:"openId"`
	Score     float32 `json:"score"`
	LeftScore float32 `json:"leftScore"`
	Reason    string  `json:"reason"`
	UserName  string  `json:"userName"`
	First     string  `json:"first,omitempty"`
	Time      int64   `json:"time,omitempty"`
	Remark    string  `json:"remark,omitempty"`
}

type WechatResponse struct {
	Code int64       `json:"state"`
	Msg  string      `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

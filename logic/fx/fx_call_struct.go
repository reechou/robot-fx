package fx

const (
	RspCodeOK = iota
	RspCodeErr
)

type CreateFxAccountReq struct {
	UnionId   string `json:"unionId"`
	WXAccount string `json:"wxAccount"`
	OpenId    string `json:"openId"`
	Name      string `json:"name"`
	Superior  string `json:"superiorId"`
}

type CreateSalesmanReq struct {
	UnionId string `json:"unionId"`
	Ticket  string `json:"ticket"`
	Phone   string `json:"phone"`
}

type updateFxBaseInfoReq struct {
	WXAccount string `json:"wxAccount"`
	OpenId    string `json:"openId"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
}

type updateFxStatusReq struct {
	UnionId string `json:"unionId"`
	Status  int64  `json:"status"`
}

type updateFxSignTimeReq struct {
	UnionId string `json:"unionId"`
}

type getFxAccountReq struct {
	UnionId string `json:"unionId"`
}

type getFxAccountFollowReq struct {
	WXAccount string `json:"wxAccount"`
	OpenId    string `json:"openId"`
}

type getFxLowerPeopleListReq struct {
	UnionId string `json:"unionId"`
	Offset  int64  `json:"offset"`
	Num     int64  `json:"num"`
}

type getFxAccountRankReq struct {
	Offset int64 `json:"offset"`
	Num    int64 `json:"num"`
}

type getFxAccountHistoryListReq struct {
	UnionId string `json:"unionId"`
	Offset  int64  `json:"offset"`
	Num     int64  `json:"num"`
}

type getFxAccountHistoryListByTypeReq struct {
	UnionId string `json:"unionId"`
	Offset  int64  `json:"offset"`
	Num     int64  `json:"num"`
	Type    int64  `json:"type"`
}

type createFxTeamReq struct {
	Name string `json:"name"`
}

type createFxTeamMemberReq struct {
	TeamId  int64  `json:"teamId"`
	UnionId string `json:"unionId"`
}

type getFxTeamMembersReq struct {
	FxTeamId int64 `json:"fxTeamId"`
}

type createFxOrderReq struct {
	UnionId     string  `json:"unionId"`
	OrderId     string  `json:"orderId"`
	GoodsId     string  `json:"goodsId"`
	OrderName   string  `json:"orderName"`
	Price       float32 `json:"price"`
	ReturnMoney float32 `json:"returnMoney"`
	Status      int64   `json:"status"`
}

type getFxOrderListReq struct {
	UnionId string `json:"unionId"`
	Status  int64  `json:"status"`
	Offset  int64  `json:"offset"`
	Num     int64  `json:"num"`
}

type getFxOrderSettlementRecordListReq struct {
	UnionId string `json:"unionId"`
	Offset  int64  `json:"offset"`
	Num     int64  `json:"num"`
}

type getFxOrderWaitSettlementRecordListReq struct {
	UnionId string `json:"unionId"`
	Offset  int64  `json:"offset"`
	Num     int64  `json:"num"`
}

type getFxOrderWaitSettlementSumReq struct {
	UnionId string `json:"unionId"`
}

type withdrawalMoneyReq struct {
	UnionId string  `json:"unionId"`
	Money   float32 `json:"money"`
	OpenId  string  `json:"openId"`
}

type getWithdrawalListReq struct {
	UnionId string `json:"unionId"`
	Offset  int64  `json:"offset"`
	Num     int64  `json:"num"`
	Status  int64  `json:"status"`
}

type getWithdrawalSumReq struct {
	UnionId string `json:"unionId"`
}

type getWithdrawalErrorListReq struct {
	Offset int64 `json:"offset"`
	Num    int64 `json:"num"`
}

type getWithdrawalErrorListFromNameReq struct {
	Name string `json:"name"`
}

type FxResponse struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

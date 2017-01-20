package controller

const (
	GodSalesman = "godlike"
	MaxSalesman = 99900
)

const (
	DEFAULT_MAX_WORKER   = 100
	DEFAULT_MAX_CHAN_LEN = 10000
)

const (
	FX_ORDER_SUCCESS    = 1 // 分销订单结算成功,后台结算成功
	FX_ORDER_WAIT       = 2 // 订单等待结算
	FX_ORDER_FAILED     = 3 // 订单失败
	FX_ORDER_SETTLEMENT = 4 // 淘宝结算
)

const (
	TAOBAO_ORDER_SUCCESS    = 1 // 订单成功
	TAOBAO_ORDER_PAY        = 2 // 订单付款
	TAOBAO_ORDER_INVALID    = 3 // 订单失效
	TAOBAO_ORDER_SETTLEMENT = 4 // 订单已结算
)

const (
	WITHDRAWAL_WAITING = iota
	WITHDRAWAL_DONE
)

const (
	FX_HISTORY_TYPE_SIGN = iota
	FX_HISTORY_TYPE_INVITE
	FX_HISTORY_TYPE_ORDER_0
	FX_HISTORY_TYPE_ORDER_1
	FX_HISTORY_TYPE_ORDER_2
	FX_HISTORY_TYPE_WITHDRAWAL
	FX_HISTORY_TYPE_SCORE_MALL
)

var (
	FxHistoryDescs = []string{"每日签到", "邀请下线", "订单返积分", "一级下线 %s", "二级下线 %s"}
)

const (
	DUOBB_GET_ALIMAMA_COOKIE_METHOD = "DuobbAccountService.GetAccountACFromAlimama"
)

const (
	ALIMAMA_TBK_PAYMENT_ONE_DAY = iota
	ALIMAMA_TBK_PAYMENT_30_DAY
)

const (
	ALIMAMA_TBK_PAYMENT_ONE_DAY_TIME = 10
	ALIMAMA_TBK_PAYMENT_30_DAY_TIMES = 12
)

const (
	ALIMAMA_TBK_PAYMENT_ORDER_STATUS_SUCCESS    = "订单成功"
	ALIMAMA_TBK_PAYMENT_ORDER_STATUS_PAY        = "订单付款"
	ALIMAMA_TBK_PAYMENT_ORDER_STATUS_INVALID    = "订单失效"
	ALIMAMA_TBK_PAYMENT_ORDER_STATUS_SETTLEMENT = "订单结算"
)

const (
	ALIMAMA_GET_TBK_PAYMENT = "http://pub.alimama.com/report/getTbkPaymentDetails.json?queryType=1&payStatus=&DownloadID=DOWNLOAD_REPORT_INCOME_NEW&startTime=%s&endTime=%s"
)

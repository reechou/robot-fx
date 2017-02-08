package controller

const (
	TestSalesman = "onKFit7N1MwqF-G3CTVCjS_qv6kk"
	GodSalesman  = "godlike"
	MaxSalesman  = 99900
)

const (
	DEFAULT_MAX_WORKER   = 100
	DEFAULT_MAX_CHAN_LEN = 10000
)

const (
	FX_ORDER_SUCCESS    = 1 // 订单结算成功,后台结算成功
	FX_ORDER_WAIT       = 2 // 订单等待结算
	FX_ORDER_FAILED     = 3 // 订单失败
	FX_ORDER_SETTLEMENT = 4 // 淘宝结算
)

// 提现状态
const (
	WITHDRAWAL_ALL     = 0 // 所有提现状态
	WITHDRAWAL_DONE    = 1 // 提现完成
	WITHDRAWAL_WAITING = 2 // 审核中
	WITHDRAWAL_FAIL    = 3 // 提现失败
)

// 积分历史记录类型
const (
	FX_HISTORY_TYPE_SIGN       = iota // 签到
	FX_HISTORY_TYPE_INVITE            // 邀请
	FX_HISTORY_TYPE_ORDER_0           // 订单主
	FX_HISTORY_TYPE_ORDER_1           // 1级分销
	FX_HISTORY_TYPE_ORDER_2           // 2级分销
	FX_HISTORY_TYPE_WITHDRAWAL        // 提现
	FX_HISTORY_TYPE_SCORE_MALL        // 积分商城
)

const (
	WX_WGLS_ACCOUNT      = "gh_1306ea147f00"
	WX_SEND_FIRST_ADD    = "网购猎手积分变动!"
	WX_SEND_FIRST_REMARK = "感谢您的使用！"
)

var (
	FxHistoryDescs = []string{
		"每日签到",
		"邀请下线",
		"订单返积分",
		"一级下线 %s",
		"二级下线 %s",
		"提现",
		"积分商城",
	}
)

const (
	DUOBB_GET_ALIMAMA_COOKIE_METHOD = "DuobbAccountService.GetAccountACFromAlimama"

	ROBOT_DEFAULT_GUIDE = "%s_微信群导购"

	ALIMAMA_GET_SELF_ADZONE_LIST = "http://pub.alimama.com/common/adzone/newSelfAdzone2.json?tag=29&itemId=541861440464&blockId=&t=%s&%s&pvid=10_%s_557_%s" // 时间 tb_token ip 时间
	ALIMAMA_GET_GUIDE_LIST       = "http://pub.alimama.com/common/site/generalize/guideList.json?t=%s&pvid=&%s&_input_charset=utf-8"                        // 时间 tb_token
	// data: name=232&categoryId=14&account1=ReeZhou&t=1484316609213&pvid=&_tb_token_=O3KOWnejbEq
	ALIMAMA_GUIDE_ADD = "http://pub.alimama.com/common/site/generalize/guideAdd.json"
	// data: tag=29&gcid=8&siteid=20776376&selectact=add&newadzonename=aaaa&t=1484316694651&_tb_token_=O3KOWnejbEq&pvid=10_125.119.120.94_557_1484315728932
	ALIMAMA_ADZONE_CREATE = "http://pub.alimama.com/common/adzone/selfAdzoneCreate.json"
)

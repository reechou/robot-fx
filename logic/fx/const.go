package fx

import (
	"errors"
)

var (
	REQ_ACCOUNT_GET_NONE = errors.New("cannot found req account.")
)

const (
	FX_ORDER_SUCCESS    = 1 // 订单结算成功,后台结算成功
	FX_ORDER_WAIT       = 2 // 订单等待结算
	FX_ORDER_FAILED     = 3 // 订单失败
	FX_ORDER_SETTLEMENT = 4 // 淘宝结算
)

const (
	KEYWORD_HELP                = "帮助"
	KEYWORD_USER_INFO           = "个人信息"
	KEYWORD_ORDER_INFO          = "订单查询"
	KEYWORD_WITHDRAWAL          = "提现"
	KEYWORD_BIND_WECHAT         = "绑定微信号"
	KEYWORD_SIGN                = "签到"
	KEYWORD_LOWER_PEOPLE        = "下线查询"
	KEYWORD_GOODS_SEARCH_URL    = "http"
	KEYWORD_GOODS_SEARCH_QUERY1 = "买"
	KEYWORD_GOODS_SEARCH_QUERY2 = "找"
	KEYWORD_GOODS_SEARCH_NEXT   = "下一个"
	KEYWORD_GOODS_SEARCH_NEXT2  = "换一个"

	KEYWORD_USER_INFO_ID         = "1"
	KEYWORD_SIGN_ID              = "2"
	KEYWORD_LOWER_PEOPLE_ID      = "3"
	KEYWORD_ORDER_INFO_ID        = "4"
	KEYWORD_WITHDRAWAL_ID        = "5"
	KEYWORD_HELP_ID              = "6"
	KEYWORD_GOODS_SEARCH_NEXT_ID = "n"
)

const (
	DEFAULT_RETURN_RATE = 50.0
	GodRate             = 0.85
)

const (
	CALLBACK_HELP = "帮助:\n* 回复 '个人信息' 或 '1' :\n查询个人信息，积分、订单总数等\n* 回复 '签到' 或 '2' :\n签到获取积分,每天只能签到一次\n" +
		"* 回复 '下线查询' 或 '3' :\n查询下线的总人数以及下线列表\n* 回复 '订单查询' 或 '4' :\n查询最近的订单记录以及具体信息\n" +
		"* 回复 '提现' 或 '5' :\n将所有可提现的积分兑换为等值的金额到您的账户上\n" +
		"* 回复 '买 xxx' 或 '找 xxx' :\n根据关键字找商品,例如: 找毛衣\n" +
		"* 回复 '下一个' 或 'n' :\n若找到的商品不满意,可发送寻找下一个商品\n" +
		"\n点击下方链接可了解如何下单得返现和下线积分规则. http://t.cn/RJwZCu0"
	CALLBACK_CREATE_ACCOUNT_SUCCESS        = "创建账户成功."
	CALLBACK_CREATE_ACCOUNT_WITHOUT_WECHAT = "创建账户成功,但未绑定微信号. 不绑定微信,会导致订单返积分不成功.\n绑定命令(范例): 绑定微信号xxx"
	CALLBACK_INVITE_SUCCESS                = "邀请 [ %s ] 成功,增加积分 [ %d ] 分."
	CALLBACK_SIGN_SUCCESS                  = "Hi, %s\n签到成功,增加积分 [ %d ] 分."
	CALLBACK_SIGN_FAILED                   = "Hi, %s\n今天已经签过到了哦."
	CALLBACK_USER_INFO_SUCCESS             = "Hi, %s\n可提现积分: %d\n历史总积分: %d\n已提现积分: %d\n订单总数: %d\n待确认订单: %d\n待确认积分: %d\n\n* 注: 100 积分 = 1 元"
	CALLBACK_LOWER_PEOPLE_SUCCESS          = "下线总人数: %d\n下线列表(只显示最近20人):"
	CALLBACK_BIND_WECHAT_SUCCESS           = "[ %s ] 绑定微信号 [ %s ] 成功."
	CALLBACK_BIND_WECHAT_FAILED            = "[ %s ] 已绑定微信号."
	CALLBACK_ORDER_LIST_SUCCESS            = "最近订单记录如下(只显示最近10条):"
	CALLBACK_GOODS_SEARCH_SUCCESS          = "Hi, %s\n【商品名称】\n%s\n★ [原价] %.02f 元\n★ [总优惠后价格约] %.02f 元\n" +
		"【优惠详情】\n★ [总优惠] %.02f 元\n★ [优惠券] %.02f 元\n★ [返利约] %.02f 元\n\n* 确认收货好评后,即可返现到可提现积分"
	CALLBACK_GOODS_SEARCH_NO_QUAN_SUCCESS = "Hi, %s\n【商品名称】\n%s\n★ [原价] %.02f 元\n★ [总优惠后价格约] %.02f 元\n" +
		"【优惠详情】\n★ [返利约] %.02f 元\n\n* 确认收货好评后,即可返现到可提现积分"
	CALLBACK_PLACE_ORDER        = "【***下单***】\n%s 长按复制本条信息,打开【手机淘宝】即可购买下单"
	CALLBACK_QUERY_PLACE_ORDER  = "【***下单***】\n%s 长按复制本条信息,打开【手机淘宝】即可购买下单\n\n* 若以上商品不满意,可回复'下一个'或'n'查找其他同类商品"
	CALLBACK_GOODS_NO_DISCOUNT  = "hi, %s\n该商品没有优惠哦"
	CALLBACK_QUERY_NO_DISCOUNT  = "hi, %s\n未找到优惠商品哦"
	CALLBACK_WITHDRAWAL_SUCCESS = "Hi, %s\n恭喜您,成功申请提现 %d 积分, 约 %.02f 元\n客服会在24小时内发放微信红包,请耐心等待.\n" +
		"每个月最多提现 %d 次.\n每次提现最少提现 %d 积分."
	CALLBACK_WITHDRAWAL_POLICY = "每个月最多提现 %d 次.\n每次提现最少提现 %d 积分."
	CALLBACK_WITHDRAWAL_NOTIFY = "返利机器人[ %s ]\n***申请提现操作***\n用户: %s\n微信号: %s\n提现金额: %.02f 元\n操作时间: %s"
)

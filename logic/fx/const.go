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
	KEYWORD_USER_INFO        = "个人信息"
	KEYWORD_ORDER_INFO       = "订单查询"
	KEYWORD_WITHDRAWAL       = "提现"
	KEYWORD_BIND_WECHAT      = "绑定微信号"
	KEYWORD_SIGN             = "签到"
	KEYWORD_LOWER_PEOPLE     = "下线查询"
	KEYWORD_GOODS_SEARCH_URL = "http"
)

const (
	DEFAULT_RETURN_RATE = 50.0
)

const (
	CALLBACK_CREATE_ACCOUNT_SUCCESS        = "创建账户成功."
	CALLBACK_CREATE_ACCOUNT_WITHOUT_WECHAT = "创建账户成功,但未绑定微信号. 不绑定微信,会导致订单返积分不成功.\n绑定命令(范例): 绑定微信号xxx"
	CALLBACK_INVITE_SUCCESS                = "邀请 [ %s ] 成功,增加积分 [ %d ] 分."
	CALLBACK_SIGN_SUCCESS                  = "签到成功,增加积分 [ %d ] 分."
	CALLBACK_SIGN_FAILED                   = "今天已经签过到了哦."
	CALLBACK_USER_INFO_SUCCESS             = "%s\n可提现积分: %d\n历史总积分: %d\n已提现积分: %d\n订单总数: %d\n待确认订单: %d\n待确认积分: %d\n\n* 注: 100 积分 = 1 元"
	CALLBACK_LOWER_PEOPLE_SUCCESS          = "下线总人数: %d\n下线列表(只显示最近20人):"
	CALLBACK_BIND_WECHAT_SUCCESS           = "[ %s ] 绑定微信号 [ %s ] 成功."
	CALLBACK_BIND_WECHAT_FAILED            = "[ %s ] 已绑定微信号."
	CALLBACK_ORDER_LIST_SUCCESS            = "最近订单记录如下(只显示最近10条):"
	CALLBACK_GOODS_SEARCH_SUCCESS          = "Hi, %s\n【商品名称】\n%s\n[原价] %.02f 元\n[总优惠后价格约] %.02f 元\n" +
		"【优惠详情】\n[总优惠] %.02f 元\n[优惠券] %.02f 元\n[好评后优惠约] %.02f 元\n" +
		"【下单】\n%s 长按复制本条信息,打开[手机淘宝]即可领券下单"
	CALLBACK_GOODS_SEARCH_NO_QUAN_SUCCESS  = "Hi, %s\n【商品名称】\n%s\n[原价] %.02f 元\n[总优惠后价格约] %.02f 元\n" +
		"【优惠详情】\n[好评后优惠约] %.02f 元\n" +
		"【下单】\n%s 长按复制本条信息,打开[手机淘宝]即可领券下单"
	CALLBACK_WITHDRAWAL_SUCCESS            = "Hi, %s\n恭喜您,成功申请提现 %d 积分, 约 %.02f 元\n客服会在24小时内发放微信红包,请耐心等待."
)
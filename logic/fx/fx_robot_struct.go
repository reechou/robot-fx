package fx

const (
	FROM_TYPE_PEOPLE         = "people"
	FROM_TYPE_GROUP          = "group"
	RECEIVE_EVENT_MSG        = "receivemsg"
	RECEIVE_EVENT_ADD_FRIEND = "addfriend"
	CHAT_TYPE_PEOPLE         = "people"
	CHAT_TYPE_GROUP          = "group"
	MSG_TYPE_TEXT            = "text"
	MSG_TYPE_IMG             = "img"
)

const (
	UNION_ID_DELIMITER = "_$$_"
)

type CreateRobotAlimamaReq struct {
	RobotWX string `json:"robotWx"`
	Alimama string `json:"alimama"`
}

type BaseInfo struct {
	Uin           string `json:"uin"`
	WechatNick    string `json:"wechatNick,omitempty"` // 微信昵称
	ReceiveEvent  string `json:"receiveEvent,omitempty"`
	FromType      string `json:"fromType,omitempty"`
	FromUserName  string `json:"fromUserName,omitempty"`
	FromNickName  string `json:"fromNickName,omitempty"`
	FromGroupName string `json:"fromGroupName,omitempty"`
}

type SendBaseInfo struct {
	WechatNick string `json:"wechatNick,omitempty"` // 微信昵称
	ChatType   string `json:"chatType,omitempty"`
	NickName   string `json:"nickName,omitempty"`
	UserName   string `json:"userName,omitempty"`
	MsgType    string `json:"msgType,omitempty"`
	Msg        string `json:"msg,omitempty"`
}

type RetResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

type AddFriend struct {
	SourceWechat string `json:"sourceWechat,omitempty"`
	SourceNick   string `json:"sourceNick,omitempty"`
	UserWxid     string `json:"userWxid,omitempty"`
	UserWechat   string `json:"userWechat,omitempty"`
	UserNick     string `json:"userNick,omitempty"`
	Ticket       string `json:"-"` // for verify
}

type ReceiveMsgInfo struct {
	BaseInfo  `json:"baseInfo,omitempty"`
	AddFriend `json:"addFriend,omitempty"`

	Msg string `json:"msg,omitempty"`
}

type CallbackMsgInfo struct {
	RetResponse `json:"retResponse,omitempty"`
	BaseInfo    `json:"baseInfo,omitempty"`

	CallbackMsgs []SendBaseInfo `json:"msg,omitempty"`
}

type SendMsgInfo struct {
	SendMsgs []SendBaseInfo `json:"sendBaseInfo,omitempty"`
}

type SendMsgResponse struct {
	RetResponse `json:"retResponse,omitempty"`
}

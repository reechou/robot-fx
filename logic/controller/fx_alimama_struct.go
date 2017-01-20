package controller

type AlimamaOtherAdzone struct {
	MemberId int `json:"memberid"`
	SiteId int `json:"siteid"`
	Name string `json:"name"`
}

type AlimamaRspInfo struct {
	Msg string `json:"message"`
	OK bool `json:"ok"`
}

type AlimamaGetSelfAdzoneData struct {
	OtherList []AlimamaOtherAdzone `json:"otherList"`
}

type AlimamaGuideInfo struct {
	Name string `json:"name"`
	GuideId int `json:"guideID"`
	MemberId int `json:"memberID"`
}

type AlimamaGuideListData struct {
	GuideList []AlimamaGuideInfo `json:"guideList"`
}

type AlimamaAdzoneCreateData struct {
	AdzoneId int `json:"adzoneId"`
	SiteId int `json:"siteId"`
}

type AlimamaGetSelfAdzoneRsp struct {
	Data AlimamaGetSelfAdzoneData `json:"data"`
	Info AlimamaRspInfo `json:"info"`
	OK bool `json:"ok"`
}

type AlimamaGetGuideListRsp struct {
	Data AlimamaGuideListData `json:"data"`
	Info AlimamaRspInfo `json:"info"`
	OK bool `json:"ok"`
}

type AlimamaAddGuideRsp struct {
	Info AlimamaRspInfo `json:"info"`
	OK bool `json:"ok"`
}

type AlimamaAdzoneCreateRsp struct {
	Data AlimamaAdzoneCreateData `json:"data"`
	Info AlimamaRspInfo `json:"info"`
	OK bool `json:"ok"`
}

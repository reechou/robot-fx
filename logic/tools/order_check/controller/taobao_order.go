package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/extrame/xls"
	"github.com/gorilla/rpc/json"
	"github.com/reechou/duobb_proto"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
	"github.com/reechou/robot-fx/logic/tools/order_check/fx_models"
)

type TaobaoOrderCheck struct {
	cfg    *config.Config
	client *http.Client
	fom    *FxOrderManager

	oneDayTimes int

	stop chan struct{}
	done chan struct{}
}

func NewTaobaoOrderCheck(cfg *config.Config, fom *FxOrderManager) *TaobaoOrderCheck {
	toc := &TaobaoOrderCheck{
		cfg:    cfg,
		client: &http.Client{},
		fom:    fom,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	//go toc.run()

	return toc
}

func (self *TaobaoOrderCheck) TaobaoOrder(reqOrder *fx_models.TaobaoOrder) error {
	if reqOrder.OrderState == TAOBAO_ORDER_INVALID {
		return nil
	}

	account := &fx_models.FxAccount{
		WechatUnionId: reqOrder.AdName,
	}
	has, err := fx_models.GetFxAccountFromWxUnionId(account)
	if err != nil {
		logrus.Errorf("get fx account error: %v", err)
		return err
	}
	if !has {
		return nil
	}
	logrus.Debugf("valid order: %v", reqOrder)

	order := &fx_models.FxOrder{
		OrderId: reqOrder.OrderId,
		GoodsId: reqOrder.GoodsId,
		Price:   reqOrder.PayPrice,
	}
	has, err = fx_models.GetFxOrderInfo(order)
	if err != nil {
		logrus.Errorf("get fx order error: %v", err)
		return err
	}
	// create order
	if !has {
		fxWxAccount := &fx_models.FxWxAccount{
			WxId: account.WxId,
		}
		wxHas, err := fx_models.GetFxWxAccount(fxWxAccount)
		if err != nil {
			logrus.Errorf("get fx wx account: %v", err)
		} else {
			if wxHas {
				order.AccountId = fxWxAccount.ID
				order.UnionId = fxWxAccount.WxId
				order.OrderName = reqOrder.GoodsInfo
				order.ReturnMoney = reqOrder.PredictingEffect
				order.AdName = reqOrder.AdName
				err = self.fom.CreateFxOrder(order)
				if err != nil {
					logrus.Errorf("create order error: %v", err)
				}
			}
		}
	}
	// update or insert tao order
	tOrder := &fx_models.TaobaoOrder{
		GoodsId:  reqOrder.GoodsId,
		PayPrice: reqOrder.PayPrice,
		OrderId:  reqOrder.OrderId,
	}
	has, err = fx_models.GetTaobaoOrder(tOrder)
	if err != nil {
		logrus.Errorf("get tao order error: %v", err)
		return err
	}
	if has {
		if tOrder.OrderState != reqOrder.OrderState {
			tOrder.OrderState = reqOrder.OrderState
			err = fx_models.UpdateTaobaoOrderStatus(tOrder)
			if err != nil {
				logrus.Errorf("update taobao order status error: %v", err)
				return err
			}
		}
		return nil
	}
	err = fx_models.CreateTaobaoOrder(reqOrder)
	if err != nil {
		logrus.Errorf("create taobao order error: %v", err)
		return err
	}

	return nil
}

func (self *TaobaoOrderCheck) run() {
	logrus.Debugf("taobao order check run start.")
	self.runCheck()
	for {
		select {
		case <-time.After(ALIMAMA_TBK_PAYMENT_ONE_DAY_TIME * time.Minute):
			self.runCheck()
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *TaobaoOrderCheck) runCheck() {
	t := ALIMAMA_TBK_PAYMENT_ONE_DAY
	if self.oneDayTimes == ALIMAMA_TBK_PAYMENT_30_DAY_TIMES {
		t = ALIMAMA_TBK_PAYMENT_30_DAY
		self.oneDayTimes = 0
	}
	self.oneDayTimes++

	alimamaList, err := fx_models.GetFxRobotAlimamaList()
	if err != nil {
		logrus.Errorf("get fx robot alimama list error: %v", err)
		return
	}
	hasChecked := make(map[string]int)
	for _, v := range alimamaList {
		if _, ok := hasChecked[v.Alimama]; ok {
			continue
		}
		self.ParseAlimamaExcel(v.Alimama, t)
		hasChecked[v.Alimama] = 1
	}
}

func (self *TaobaoOrderCheck) ParseAlimamaExcel(ali string, t int) {
	// get alimama cookie
	cookie := self.getAlimamaCookie(ali)
	if cookie == "" {
		logrus.Errorf("cannot get the cookie.")
		return
	}
	logrus.Debug(cookie)

	// get alimama excel
	var start, end string
	end = time.Now().Format("2006-01-02")
	if t == ALIMAMA_TBK_PAYMENT_ONE_DAY {
		// 1 day
		start = end
	} else {
		start = time.Unix(time.Now().Unix()-30*86400, 0).Format("2006-01-02")
	}

	// get excel
	body, err := self.getAlimamaTbkPaymentExcel(fmt.Sprintf(ALIMAMA_GET_TBK_PAYMENT, start, end), cookie)
	if err != nil {
		logrus.Errorf("get excel error: %v", err)
		return
	}

	// parse excel
	xlFile, err := xls.OpenReader(bytes.NewReader(body), "utf-8")
	if err != nil {
		logrus.Errorf("excel open reader error: %v", err)
		return
	}
	res := xlFile.ReadAllCells(int(xlFile.GetSheet(0).MaxRow))
	for _, v := range res {
		if len(v) != 30 {
			logrus.Errorf("len(v)[%v] != 30", v)
			continue
		}
		//fmt.Println(len(v), v)
		if v[8] == ALIMAMA_TBK_PAYMENT_ORDER_STATUS_INVALID {
			continue
		}
		account := &fx_models.FxAccount{
			WechatUnionId: v[29],
		}
		has, err := fx_models.GetFxAccountFromWxUnionId(account)
		if err != nil {
			logrus.Errorf("get fx account error: %v", err)
			continue
		}
		if !has {
			continue
		}
		logrus.Debugf("valid order: %v", v)

		payMoney, err := strconv.ParseFloat(v[12], 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		returnMoney, err := strconv.ParseFloat(v[13], 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}

		order := &fx_models.FxOrder{
			OrderId: v[24],
			GoodsId: v[3],
			Price:   float32(payMoney),
		}
		has, err = fx_models.GetFxOrderInfo(order)
		if err != nil {
			logrus.Errorf("get fx order error: %v", err)
			continue
		}
		// create order
		if !has {
			order.AccountId = account.ID
			order.UnionId = account.WxId
			order.OrderName = v[2]
			order.ReturnMoney = float32(returnMoney)
			err = self.fom.CreateFxOrder(order)
			if err != nil {
				logrus.Errorf("create order error: %v", err)
			}
		}
		// update or insert tao order
		orderState := TAOBAO_ORDER_SUCCESS
		if v[8] == ALIMAMA_TBK_PAYMENT_ORDER_STATUS_PAY {
			orderState = TAOBAO_ORDER_PAY
		} else if v[8] == ALIMAMA_TBK_PAYMENT_ORDER_STATUS_SETTLEMENT {
			orderState = TAOBAO_ORDER_SETTLEMENT
		}
		tOrder := &fx_models.TaobaoOrder{
			GoodsId:  v[3],
			PayPrice: float32(payMoney),
			OrderId:  v[24],
		}
		has, err = fx_models.GetTaobaoOrder(tOrder)
		if err != nil {
			logrus.Errorf("get tao order error: %v", err)
			continue
		}
		if has {
			if tOrder.OrderState != orderState {
				tOrder.OrderState = orderState
				err = fx_models.UpdateTaobaoOrderStatus(tOrder)
				if err != nil {
					logrus.Errorf("update taobao order status error: %v", err)
				}
			}
			continue
		}
		goodsNum, err := strconv.Atoi(v[6])
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		goodsPrice, err := strconv.ParseFloat(v[7], 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		incomeRatio, err := strconv.ParseFloat(self.filterString(v[10]), 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		splitRatio, err := strconv.ParseFloat(self.filterString(v[11]), 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		predictingEffect, err := strconv.ParseFloat(v[13], 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		settlementMoney, err := strconv.ParseFloat(v[14], 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		estimatedIncome, err := strconv.ParseFloat(v[15], 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		commissionRate, err := strconv.ParseFloat(self.filterString(v[17]), 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		commissionMoney, err := strconv.ParseFloat(v[18], 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		subsidyRate, err := strconv.ParseFloat(self.filterString(v[19]), 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		subsidyMoney, err := strconv.ParseFloat(v[20], 32)
		if err != nil {
			logrus.Errorf("strconv parse error: %v", err)
			continue
		}
		taobaoOrder := &fx_models.TaobaoOrder{
			OrderCreatedTime:    v[0],
			OrderClickTime:      v[1],
			GoodsInfo:           v[2],
			GoodsId:             v[3],
			WangwangName:        v[4],
			StoreName:           v[5],
			GoodsNum:            goodsNum,
			GoodsPrice:          float32(goodsPrice),
			OrderState:          orderState,
			OrderType:           v[9],
			IncomeRatio:         float32(incomeRatio),
			SplitRatio:          float32(splitRatio),
			PayPrice:            float32(payMoney),
			PredictingEffect:    float32(predictingEffect),
			SettlementMoney:     float32(settlementMoney),
			EstimatedIncome:     float32(estimatedIncome),
			SettlementTime:      v[16],
			CommissionRate:      float32(commissionRate),
			CommissionMoney:     float32(commissionMoney),
			SubsidyRate:         float32(subsidyRate),
			SubsidyMoney:        float32(subsidyMoney),
			SubsidyType:         v[21],
			TransactionPlatform: v[22],
			ThirdPartyService:   v[23],
			OrderId:             v[24],
			CategoryName:        v[25],
			SiteId:              v[26],
			SiteName:            v[27],
			AdId:                v[28],
			AdName:              v[29],
		}
		err = fx_models.CreateTaobaoOrder(taobaoOrder)
		if err != nil {
			logrus.Errorf("create taobao order error: %v", err)
		}
	}
}

func (self *TaobaoOrderCheck) filterString(v string) string {
	return strings.Replace(v, " %", "", -1)
}

func (self *TaobaoOrderCheck) getAlimamaCookie(ali string) string {
	request := map[string]string{"ali": ali}
	url := "http://" + self.cfg.DuobbSrv.Host + "/rpc"
	message, err := json.EncodeClientRequest(DUOBB_GET_ALIMAMA_COOKIE_METHOD, request)
	if err != nil {
		logrus.Errorf("json encode client request error: %v", err)
		return ""
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		logrus.Errorf("http new request error: %v", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := self.client.Do(req)
	if err != nil {
		logrus.Errorf("http do request error: %v", err)
		return ""
	}
	defer resp.Body.Close()
	var result duobb_proto.Response
	err = json.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		logrus.Errorf("json decode client response error: %v", err)
		return ""
	}
	if result.Code == duobb_proto.DUOBB_RSP_SUCCESS {
		dataMap := result.Data.(map[string]interface{})
		if dataMap != nil {
			cookie, ok := dataMap["cookie"]
			if ok {
				return cookie.(string)
			}
		}
	}

	return ""
}

func (self *TaobaoOrderCheck) getAlimamaTbkPaymentExcel(url, cookie string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Errorf("get excel http new request error: %v", err)
		return nil, err
	}
	req.Header.Set("Host", "pub.alimama.com")
	req.Header.Set("Referer", "http://pub.alimama.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36")
	req.Header.Set("Cookie", cookie)
	resp, err := self.client.Do(req)
	if err != nil {
		logrus.Errorf("get excel http do request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("get excel io realall error: %v", err)
		return nil, err
	}
	return body, nil
}

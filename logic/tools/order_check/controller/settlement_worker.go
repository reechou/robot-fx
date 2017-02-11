package controller

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/tools/order_check/act"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
	"github.com/reechou/robot-fx/logic/tools/order_check/ext"
	"github.com/reechou/robot-fx/logic/tools/order_check/fx_models"
)

type SettlementWorker struct {
	orderChanList []chan *fx_models.FxOrder

	cfg   *config.Config
	act   *act.ActLogic
	wrExt *ext.WxRobotExt

	wg   sync.WaitGroup
	stop chan struct{}
}

func NewSettlementWorker(maxWorker, maxChanLen int, cfg *config.Config, wrExt *ext.WxRobotExt) *SettlementWorker {
	sw := &SettlementWorker{
		cfg:   cfg,
		stop:  make(chan struct{}),
		act:   act.NewActLogic(cfg),
		wrExt: wrExt,
	}
	for i := 0; i < maxWorker; i++ {
		orderChan := make(chan *fx_models.FxOrder, maxChanLen)
		sw.orderChanList = append(sw.orderChanList, orderChan)
		sw.wg.Add(1)
		go sw.runWorker(orderChan, sw.stop)
	}
	return sw
}

func (sw *SettlementWorker) Close() {
	close(sw.stop)
	sw.wg.Wait()
}

func (sw *SettlementWorker) SettlementOrder(order *fx_models.FxOrder) {
	idx := int(order.ID) % len(sw.orderChanList)
	select {
	case sw.orderChanList[idx] <- order:
	case <-time.After(5 * time.Second):
		logrus.Errorf("settlement into order channel timeout.")
	}
}

func (sw *SettlementWorker) runWorker(orderChan chan *fx_models.FxOrder, stop chan struct{}) {
	for {
		select {
		case order := <-orderChan:
			sw.do(order)
		case <-stop:
			sw.wg.Done()
			return
		}
	}
}

func (sw *SettlementWorker) do(order *fx_models.FxOrder) {
	if order.Status != FX_ORDER_WAIT {
		logrus.Errorf("order[%v] cannot be settlement.", order)
		return
	}

	// check status
	checkOrder := &fx_models.FxOrder{
		OrderId: order.OrderId,
		GoodsId: order.GoodsId,
		Price:   order.Price,
	}
	has, err := fx_models.GetFxOrderInfo(checkOrder)
	if err != nil {
		logrus.Errorf("get fx order[%v] status error: %v", order, err)
		return
	}
	if !has {
		logrus.Errorf("get fx order[%v] has no this order", order)
		return

	}
	if checkOrder.Status != FX_ORDER_WAIT {
		logrus.Errorf("order[%v] cannot be settlement, order status: %d", order, checkOrder.Status)
		return
	}

	var levelReturns []float32
	for i := 0; i < len(sw.cfg.SettlementCommission.LevelPer); i++ {
		lReturn := order.ReturnMoney * GodRate * float32(sw.cfg.SettlementCommission.LevelPer[i]) / 100.0 * float32(sw.cfg.Score.EnlargeScale)
		levelReturns = append(levelReturns, lReturn)
	}

	settlementFxOrder := &fx_models.SettlementFxOrderInfo{
		Status:        FX_ORDER_SUCCESS,
		Order:         order,
		OrderAddMoney: levelReturns[0],
	}
	err = fx_models.SettlementOwnerFxOrder(settlementFxOrder)
	if err != nil {
		logrus.Errorf("do settlement order[%v] settlement owner order error: %v", order, err)
		return
	}
	logrus.Infof("order_id[%s][%s][%f] settlement for owner[%s] with return_money[%f] success", order.OrderId, order.GoodsId, order.Price, order.UnionId, levelReturns[0])

	now := time.Now().Unix()

	//month := fmt.Sprintf(time.Now().Format("200601"))
	//err = sw.updateFxAccountMonth(month, order.UnionId, levelReturns[0])
	//if err != nil {
	//	logrus.Errorf("do settlement order[%v] update fx account month owner order error: %v", order, err)
	//	return err
	//}

	orderFxWxAccount := &fx_models.FxWxAccount{
		WxId: order.UnionId,
	}
	has, err = fx_models.GetFxWxAccount(orderFxWxAccount)
	if err != nil {
		logrus.Errorf("do settlement order[%v] in level[0] get fx account from wx_id[%s] error: %v",
			order, order.UnionId, err)
		return
	}
	if !has {
		logrus.Errorf("do settlement no this owner account[%s]", order.UnionId)
		return
	}

	var recordList []fx_models.FxOrderSettlementRecord
	recordList = append(recordList, fx_models.FxOrderSettlementRecord{
		AccountId:    orderFxWxAccount.ID,
		UnionId:      order.UnionId,
		OrderId:      order.OrderId,
		GoodsId:      order.GoodsId,
		Price:        order.Price,
		ReturnMoney:  levelReturns[0],
		SourceId:     order.UnionId,
		Level:        0,
		OrderCreated: order.CreatedAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	})

	var historyList []fx_models.FxAccountHistory
	historyList = append(historyList, fx_models.FxAccountHistory{
		AccountId:  orderFxWxAccount.ID,
		UnionId:    orderFxWxAccount.WxId,
		Score:      levelReturns[0],
		ChangeType: int64(FX_HISTORY_TYPE_ORDER_0),
		ChangeDesc: FxHistoryDescs[FX_HISTORY_TYPE_ORDER_0],
		CreatedAt:  now,
	})

	var robotWx string
	adList := strings.Split(order.AdName, ext.UNION_ID_DELIMITER)
	if len(adList) == 2 {
		robotWx = adList[0]
	}
	orderOwnerName := orderFxWxAccount.Name
	var notifyMsgs ext.SendMsgInfo
	notifyMsgs.SendMsgs = append(notifyMsgs.SendMsgs, ext.SendBaseInfo{
		WechatNick: robotWx,
		ChatType:   ext.CHAT_TYPE_PEOPLE,
		NickName:   orderOwnerName,
		MsgType:    ext.MSG_TYPE_TEXT,
		Msg:        fmt.Sprintf(NOTIFY_MSG_SETTLEMENT_ORDER_OWNER, order.OrderId[:4], int64(levelReturns[0])),
	})

	var upperFxAccount *fx_models.FxWxAccount

	// this unionId is always wx_id
	unionId := orderFxWxAccount.Superior
	for i := 1; i < len(levelReturns); i++ {
		// get upper
		fxWxAccount := &fx_models.FxWxAccount{
			WxId: unionId,
		}
		has, err := fx_models.GetFxWxAccount(fxWxAccount)
		if err != nil {
			logrus.Errorf("do settlement order[%v] in level[%d] get fx wx account from wx_id[%s] error: %v",
				order, i, unionId, err)
			return
		}
		if !has {
			logrus.Debugf("do settlement no this account of wx_id[%s]", unionId)
			break
		}
		if i == 1 {
			upperFxAccount = &fx_models.FxWxAccount{
				ID:   fxWxAccount.ID,
				WxId: fxWxAccount.WxId,
				Name: fxWxAccount.Name,
			}
		}
		// add return money
		err = fx_models.AddFxWxAccountMoney(levelReturns[i], fxWxAccount)
		if err != nil {
			logrus.Errorf("do settlement order[%v] in level[%d] add money in fx account from union_id[%d] error: %v",
				order, i, unionId, err)
			return
		}
		logrus.Infof("order_id[%s] settlement for upper user[%s][level-%d] with return_money[%f] success", order.OrderId, unionId, i, levelReturns[i])

		//err = sw.updateFxAccountMonth(month, unionId, levelReturns[i])
		//if err != nil {
		//	logrus.Errorf("do settlement order[%v] update fx account month union_id[%s][level-%d] order error: %v", order, unionId, i, err)
		//	return err
		//}

		//recordList = append(recordList, fx_models.FxOrderSettlementRecord{
		//	AccountId:   fxAccount.ID,
		//	UnionId:     unionId,
		//	OrderId:     order.OrderId,
		//	ReturnMoney: levelReturns[i],
		//	SourceId:    order.UnionId,
		//	Level:       int64(i),
		//	CreatedAt:   now,
		//	UpdatedAt:   now,
		//})

		historyList = append(historyList, fx_models.FxAccountHistory{
			AccountId:  fxWxAccount.ID,
			UnionId:    fxWxAccount.WxId,
			Score:      levelReturns[i],
			ChangeType: int64(FX_HISTORY_TYPE_ORDER_0 + i),
			ChangeDesc: fmt.Sprintf(FxHistoryDescs[FX_HISTORY_TYPE_ORDER_0+i], fxWxAccount.Name),
			CreatedAt:  now,
		})

		notifyMsgs.SendMsgs = append(notifyMsgs.SendMsgs, ext.SendBaseInfo{
			WechatNick: robotWx,
			ChatType:   ext.CHAT_TYPE_PEOPLE,
			NickName:   fxWxAccount.Name,
			MsgType:    ext.MSG_TYPE_TEXT,
			Msg:        fmt.Sprintf(NOTIFY_MSG_SETTLEMENT_ORDER_UPPER, i, orderOwnerName, order.OrderId[:4], int64(levelReturns[i])),
		})

		unionId = fxWxAccount.Superior
	}

	err = fx_models.CreateFxOrderSettlementRecordList(recordList)
	if err != nil {
		logrus.Errorf("create fx order[%v] settlement record list error: %v", order, err)
	}
	// insert history
	err = fx_models.CreateFxAccountHistoryList(historyList)
	if err != nil {
		logrus.Errorf("create fx order[%v] fx account history list error: %v", order, err)
	}
	err = sw.wrExt.SendMsg(&notifyMsgs)
	if err != nil {
		logrus.Errorf("send notify msg error: %v", err)
	}
	// check order act
	sw.act.CheckActOfOrder(orderFxWxAccount, upperFxAccount)
}

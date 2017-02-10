package controller

import (
	"fmt"
	"time"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
	"github.com/reechou/robot-fx/logic/tools/order_check/fx_models"
	"github.com/reechou/robot-fx/logic/tools/order_check/ext"
)

type FxOrderManager struct {
	cfg *config.Config
	wrExt *ext.WxRobotExt
}

func NewFxOrderManager(cfg *config.Config, wrExt *ext.WxRobotExt) *FxOrderManager {
	fom := &FxOrderManager{
		cfg: cfg,
		wrExt: wrExt,
	}
	return fom
}

func (self *FxOrderManager) CreateFxOrder(info *fx_models.FxOrder) error {
	info.Status = FX_ORDER_WAIT
	err := fx_models.CreateFxOrder(info)
	if err != nil {
		logrus.Errorf("create fx order error: %v", err)
		return err
	}

	if len(self.cfg.LevelPer) == 0 {
		logrus.Errorf("please check config of level percentage")
		return nil
	}

	var levelReturns []float32
	for i := 0; i < len(self.cfg.LevelPer); i++ {
		lReturn := info.ReturnMoney * GodRate * float32(self.cfg.LevelPer[i]) / 100.0 * float32(self.cfg.Score.EnlargeScale)
		levelReturns = append(levelReturns, lReturn)
	}

	var notifyMsgs ext.SendMsgInfo
	var robotWx string
	adList := strings.Split(info.AdName, ext.UNION_ID_DELIMITER)
	if len(adList) == 2 {
		robotWx = adList[0]
	}
	
	now := time.Now().Unix()
	var recordList []fx_models.FxOrderWaitSettlementRecord
	recordList = append(recordList, fx_models.FxOrderWaitSettlementRecord{
		AccountId:   info.AccountId,
		UnionId:     info.UnionId,
		OrderId:     info.OrderId,
		GoodsId:     info.GoodsId,
		Price:       info.Price,
		ReturnMoney: levelReturns[0],
		Level:       0,
		CreatedAt:   now,
	})

	fxWxAccount := &fx_models.FxWxAccount{
		WxId: info.UnionId,
	}
	has, err := fx_models.GetFxWxAccount(fxWxAccount)
	if err != nil {
		logrus.Errorf("create order[%v] get fx wx account from wx_id[%s] error: %v",
			info, info.UnionId, err)
		return err
	}
	if !has {
		logrus.Errorf("create order no this owern account wx_id[%s]", info.UnionId)
		return fmt.Errorf("create order no this owern account wx_id[%s]", info.UnionId)
	}
	
	notifyMsgs.SendMsgs = append(notifyMsgs.SendMsgs, ext.SendBaseInfo{
		WechatNick: robotWx,
		ChatType:   ext.CHAT_TYPE_PEOPLE,
		NickName:   fxWxAccount.Name,
		MsgType:    ext.MSG_TYPE_TEXT,
		Msg:        fmt.Sprintf(NOTIFY_MSG_CREATE_ORDER_OWNER, info.OrderId[:4], int64(levelReturns[0])),
	})

	unionId := fxWxAccount.Superior
	for i := 1; i < len(levelReturns); i++ {
		// get upper
		fxWxAccount := &fx_models.FxWxAccount{
			WxId: unionId,
		}
		if fxWxAccount.WxId == GodSalesman {
			break
		}
		has, err := fx_models.GetFxWxAccount(fxWxAccount)
		if err != nil {
			logrus.Errorf("create wait settlement order[%v] in level[%d] get fx wx account from wx_id[%s] error: %v",
				info, i, unionId, err)
			return err
		}
		if !has {
			logrus.Debugf("create wait settlement no this account wx_id[%s]", unionId)
			break
		}

		recordList = append(recordList, fx_models.FxOrderWaitSettlementRecord{
			AccountId:   fxWxAccount.ID,
			UnionId:     unionId,
			OrderId:     info.OrderId,
			GoodsId:     info.GoodsId,
			Price:       info.Price,
			ReturnMoney: levelReturns[i],
			Level:       int64(i),
			CreatedAt:   now,
		})
		
		notifyMsgs.SendMsgs = append(notifyMsgs.SendMsgs, ext.SendBaseInfo{
			WechatNick: robotWx,
			ChatType:   ext.CHAT_TYPE_PEOPLE,
			NickName:   fxWxAccount.Name,
			MsgType:    ext.MSG_TYPE_TEXT,
			Msg:        fmt.Sprintf(NOTIFY_MSG_CREATE_ORDER_UPPER, i, fxWxAccount.Name, info.OrderId[:4], int64(levelReturns[i])),
		})

		unionId = fxWxAccount.Superior
	}

	err = fx_models.CreateFxOrderWaitSettlementRecordList(recordList)
	if err != nil {
		logrus.Errorf("create fx order[%d] wait settlement record list error: %v", info, err)
		return err
	}
	
	err = self.wrExt.SendMsg(notifyMsgs)
	if err != nil {
		logrus.Errorf("send notify msg error: %v", err)
	}

	return nil
}

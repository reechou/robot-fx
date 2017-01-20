package controller

import (
	"fmt"
	"time"
	
	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
	"github.com/reechou/robot-fx/logic/tools/order_check/fx_models"
)

type FxOrderManager struct {
	cfg *config.Config
}

func NewFxOrderManager(cfg *config.Config) *FxOrderManager {
	fom := &FxOrderManager{
		cfg: cfg,
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
		lReturn := info.ReturnMoney * float32(self.cfg.LevelPer[i]) / 100.0 * float32(self.cfg.Score.EnlargeScale)
		levelReturns = append(levelReturns, lReturn)
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
	
	fxAccount := &fx_models.FxAccount{
		UnionId: info.UnionId,
	}
	has, err := fx_models.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("create order[%v] get fx account from union_id[%d] error: %v",
			info, info.UnionId, err)
		return err
	}
	if !has {
		logrus.Errorf("create order no this owern account[%s]", info.UnionId)
		return fmt.Errorf("create order no this owern account[%s]", info.UnionId)
	}
	
	unionId := fxAccount.Superior
	for i := 1; i < len(levelReturns); i++ {
		// get upper
		fxAccount := &fx_models.FxAccount{
			UnionId: unionId,
		}
		if fxAccount.UnionId == GodSalesman {
			break
		}
		has, err := fx_models.GetFxAccount(fxAccount)
		if err != nil {
			logrus.Errorf("create wait settlement order[%v] in level[%d] get fx account from union_id[%d] error: %v",
				info, i, unionId, err)
			return err
		}
		if !has {
			logrus.Debugf("create wait settlement no this account[%s]", unionId)
			break
		}
		
		recordList = append(recordList, fx_models.FxOrderWaitSettlementRecord{
			AccountId:   fxAccount.ID,
			UnionId:     unionId,
			OrderId:     info.OrderId,
			GoodsId:     info.GoodsId,
			Price:       info.Price,
			ReturnMoney: levelReturns[i],
			Level:       int64(i),
			CreatedAt:   now,
		})
		
		unionId = fxAccount.Superior
	}
	
	err = fx_models.CreateFxOrderWaitSettlementRecordList(recordList)
	if err != nil {
		logrus.Errorf("create fx order[%d] wait settlement record list error: %v", info, err)
		return err
	}
	
	return nil
}

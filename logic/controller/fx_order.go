package controller

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/models"
)

func (daemon *Daemon) CreateFxOrder(info *models.FxOrder) error {
	info.Status = FX_ORDER_WAIT
	err := models.CreateFxOrder(info)
	if err != nil {
		logrus.Errorf("create fx order error: %v", err)
		return err
	}

	if len(daemon.cfg.LevelPer) == 0 {
		logrus.Errorf("please check config of level percentage")
		return nil
	}

	var levelReturns []float32
	for i := 0; i < len(daemon.cfg.LevelPer); i++ {
		lReturn := info.ReturnMoney * float32(daemon.cfg.LevelPer[i]) / 100.0 * float32(daemon.cfg.Score.EnlargeScale)
		levelReturns = append(levelReturns, lReturn)
	}

	now := time.Now().Unix()
	var recordList []models.FxOrderWaitSettlementRecord
	recordList = append(recordList, models.FxOrderWaitSettlementRecord{
		AccountId:   info.AccountId,
		UnionId:     info.UnionId,
		OrderId:     info.OrderId,
		GoodsId:     info.GoodsId,
		Price:       info.Price,
		ReturnMoney: levelReturns[0],
		Level:       0,
		CreatedAt:   now,
	})

	fxAccount := &models.FxAccount{
		UnionId: info.UnionId,
	}
	has, err := models.GetFxAccount(fxAccount)
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
		fxAccount := &models.FxAccount{
			UnionId: unionId,
		}
		if fxAccount.UnionId == GodSalesman {
			break
		}
		has, err := models.GetFxAccount(fxAccount)
		if err != nil {
			logrus.Errorf("create wait settlement order[%v] in level[%d] get fx account from union_id[%d] error: %v",
				info, i, unionId, err)
			return err
		}
		if !has {
			logrus.Debugf("create wait settlement no this account[%s]", unionId)
			break
		}

		recordList = append(recordList, models.FxOrderWaitSettlementRecord{
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

	err = models.CreateFxOrderWaitSettlementRecordList(recordList)
	if err != nil {
		logrus.Errorf("create fx order[%d] wait settlement record list error: %v", info, err)
		return err
	}

	return nil
}

func (daemon *Daemon) GetFxOrderListCount(unionId string) (int64, error) {
	return models.GetFxOrderListCount(unionId)
}

func (daemon *Daemon) GetFxOrderListCountById(accountId int64) (int64, error) {
	return models.GetFxOrderListCountById(accountId)
}

func (daemon *Daemon) GetFxAllOrderList(unionId string, offset, num int64) ([]models.FxOrder, error) {
	list, err := models.GetFxAllOrderList(unionId, offset, num)
	if err != nil {
		logrus.Errorf("get fx order list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) GetFxOrderList(unionId string, offset, num, status int64) ([]models.FxOrder, error) {
	list, err := models.GetFxOrderList(unionId, offset, num, status)
	if err != nil {
		logrus.Errorf("get fx order list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) GetFxOrderListById(accountId int64, offset, num, status int64) ([]models.FxOrder, error) {
	list, err := models.GetFxOrderListById(accountId, offset, num, status)
	if err != nil {
		logrus.Errorf("get fx order list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) CreateFxOrderSettlementRecord(info *models.FxOrderSettlementRecord) error {
	err := models.CreateFxOrderSettlementRecord(info)
	if err != nil {
		logrus.Errorf("create fx order settlement record error: %v", err)
		return err
	}
	return nil
}

func (daemon *Daemon) GetFxOrderSettlementRecordListCount(unionId string) (int64, error) {
	return models.GetFxOrderSettlementRecordListCount(unionId)
}

func (daemon *Daemon) GetFxOrderSettlementRecordListCountById(accountId int64) (int64, error) {
	return models.GetFxOrderSettlementRecordListCountById(accountId)
}

func (daemon *Daemon) GetFxOrderSettlementRecordList(unionId string, offset, num int64) ([]models.FxOrderSettlementRecord, error) {
	list, err := models.GetFxOrderSettlementRecordList(unionId, offset, num)
	if err != nil {
		logrus.Errorf("get fx order settlement record list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) GetFxOrderSettlementRecordListById(accountId int64, offset, num int64) ([]models.FxOrderSettlementRecord, error) {
	list, err := models.GetFxOrderSettlementRecordListByid(accountId, offset, num)
	if err != nil {
		logrus.Errorf("get fx order settlement record list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) GetFxOrderWaitSettlementRecordSum(accountId int64) (float32, error) {
	return models.GetFxOrderWaitSettlementRecordListSumById(accountId, FX_ORDER_WAIT)
}

func (daemon *Daemon) GetFxOrderWaitSettlementRecordListCount(unionId string) (int64, error) {
	return models.GetFxOrderWaitSettlementRecordListCount(unionId, FX_ORDER_WAIT)
}

func (daemon *Daemon) GetFxOrderWaitSettlementRecordListCountById(accountId int64) (int64, error) {
	return models.GetFxOrderWaitSettlementRecordListCountById(accountId, FX_ORDER_WAIT)
}

func (daemon *Daemon) GetFxOrderWaitSettlementRecordListById(accountId int64, offset, num int64) ([]models.FxOrderWaitSettlementRecord, error) {
	list, err := models.GetFxOrderWaitSettlementRecordListById(accountId, offset, num, FX_ORDER_WAIT)
	if err != nil {
		logrus.Errorf("get fx wait settlement order record list error: %v", err)
		return nil, err
	}
	return list, err
}

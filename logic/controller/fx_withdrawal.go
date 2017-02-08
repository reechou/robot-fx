package controller

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/models"
)

// return error, ifSystemError
func (daemon *Daemon) CreateWithdrawalRecord(info *models.WithdrawalRecord, fxAccount *models.FxAccount, fxWxAccount *models.FxWxAccount) (error, bool) {
	if info.WithdrawalMoney < float32(daemon.cfg.WithdrawalPolicy.MinimumWithdrawal) {
		return ErrWithdrawalMinimum, false
	}

	if fxWxAccount.CanWithdrawals < info.WithdrawalMoney {
		return ErrWithdrawalLimitBalance, false
	}

	monthCount, err := models.GetMonthWithdrawalRecord(info.UnionId)
	if err != nil {
		logrus.Errorf("get month withdrawal record error: %v", err)
		return err, true
	}
	if monthCount >= int64(daemon.cfg.WithdrawalPolicy.MonthWithdrawalTime) {
		return ErrWithdrawalOverMonthLimit, false
	}

	if daemon.cfg.WithdrawalPolicy.IfWithdrawalCheck {
		err = models.MinusFxWxAccountMoney(info.WithdrawalMoney, fxWxAccount)
		if err != nil {
			logrus.Errorf("withdrawal money[%f] with account[%v] error: %v", info.WithdrawalMoney, fxWxAccount, err)
			return err, true
		}

		info.AccountId = fxWxAccount.ID
		info.RobotWx = fxAccount.RobotWx
		info.Wechat = fxAccount.Wechat
		info.Name = fxAccount.Name
		info.Status = WITHDRAWAL_WAITING
		info.Balance = fxWxAccount.CanWithdrawals - info.WithdrawalMoney
		err = models.CreateWithdrawalRecord(info)
		if err != nil {
			logrus.Errorf("create withdrawal record error: %v", err)
			return err, true
		}

		h := models.FxAccountHistory{
			UnionId:    info.UnionId,
			Score:      -info.WithdrawalMoney,
			ChangeType: int64(FX_HISTORY_TYPE_WITHDRAWAL),
			ChangeDesc: FxHistoryDescs[FX_HISTORY_TYPE_WITHDRAWAL],
			CreatedAt:  time.Now().Unix(),
		}
		models.CreateFxAccountHistoryList([]models.FxAccountHistory{h})

		// 直接提现
		//wReq := &ext.WithdrawalReq{
		//	OpenId:      info.OpenId,
		//	TotalAmount: int64(info.WithdrawalMoney / float32(daemon.cfg.Score.EnlargeScale) * 100),
		//	MchBillno:   strconv.Itoa(int(info.ID)),
		//}
		//err = daemon.we.Withdrawal(wReq)
		//if err != nil {
		//	logrus.Errorf("account[%v] info[%v] wechat withdrawal error: %v", fxAccount, info, err)
		//	wErrInfo := &models.WithdrawalRecordError{
		//		AccountId:       fxAccount.ID,
		//		UnionId:         fxAccount.UnionId,
		//		Name:            fxAccount.Name,
		//		WithdrawalMoney: info.WithdrawalMoney,
		//		ErrorMsg:        err.Error(),
		//	}
		//	err = models.CreateWithdrawalRecordError(wErrInfo)
		//	if err != nil {
		//		logrus.Errorf("info[%v] create withdrawal error msg record error: %v", wErrInfo, err)
		//	}
		//} else {
		//	logrus.Infof("user[%s] withdrawl[%f] wechat success.", info.UnionId)
		//}
	} else {
		info.AccountId = fxAccount.ID
		info.Status = WITHDRAWAL_WAITING
		err = models.CreateWithdrawalRecord(info)
		if err != nil {
			logrus.Errorf("create withdrawal record error: %v", err)
			return err, true
		}
	}

	return nil, false
}

func (daemon *Daemon) GetWithdrawalRecordListCount(unionId string, status int64) (int64, error) {
	return models.GetWithdrawalRecordListCount(unionId, status)
}

func (daemon *Daemon) GetWithdrawalRecordListCountById(accountId int64) (int64, error) {
	return models.GetWithdrawalRecordListCountById(accountId)
}

func (daemon *Daemon) GetWithdrawalRecordList(unionId string, offset, num, status int64) ([]models.WithdrawalRecord, error) {
	list, err := models.GetWithdrawalRecordList(unionId, offset, num, status)
	if err != nil {
		logrus.Errorf("get withdrawal record list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) GetWithdrawalRecordListById(accountId int64, offset, num int64) ([]models.WithdrawalRecord, error) {
	list, err := models.GetWithdrawalRecordListById(accountId, offset, num)
	if err != nil {
		logrus.Errorf("get withdrawal record list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) GetWithdrawalRecordSum(unionId string) (float32, error) {
	return models.GetWithdrawalRecordSum(unionId)
}

func (daemon *Daemon) GetWithdrawalErrorRecordListCount() (int64, error) {
	return models.GetWithdrawalRecordErrorListCount()
}

func (daemon *Daemon) GetWithdrawalErrorRecordList(offset, num int64) ([]models.WithdrawalRecordError, error) {
	list, err := models.GetWithdrawalRecordErrorList(offset, num)
	if err != nil {
		logrus.Errorf("get withdrawal error msg record list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) GetWithdrawalErrorRecordListFromName(name string) ([]models.WithdrawalRecordError, error) {
	list, err := models.GetWithdrawalRecordErrorListFromName(name)
	if err != nil {
		logrus.Errorf("get withdrawal error msg record list from name error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) ConfirmWithdrawal(id int64) error {
	return models.UpdateWithdrawalRecord(id, WITHDRAWAL_DONE)
}

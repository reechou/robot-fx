package controller

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/config"
	"github.com/reechou/robot-fx/logic/models"
)

type CashWithdrawalWorker struct {
	withdrawalChanList []chan *models.WithdrawalRecord

	cfg *config.Config

	wg   sync.WaitGroup
	stop chan struct{}
}

func NewCashWithdrawalWorker(maxWorker, maxChanLen int, cfg *config.Config) *CashWithdrawalWorker {
	cww := &CashWithdrawalWorker{
		cfg:  cfg,
		stop: make(chan struct{}),
	}
	for i := 0; i < maxWorker; i++ {
		wChan := make(chan *models.WithdrawalRecord, maxChanLen)
		cww.withdrawalChanList = append(cww.withdrawalChanList, wChan)
		cww.wg.Add(1)
		go cww.runWorker(wChan, cww.stop)
	}
	return cww
}

func (cww *CashWithdrawalWorker) Close() {
	close(cww.stop)
	cww.wg.Wait()
}

func (cww *CashWithdrawalWorker) Withdrawal(wInfo *models.WithdrawalRecord) {
	idx := int(wInfo.AccountId) % len(cww.withdrawalChanList)
	select {
	case cww.withdrawalChanList[idx] <- wInfo:
	case <-time.After(5 * time.Second):
		logrus.Errorf("Withdrawal into channel timeout.")
	}
}

func (cww *CashWithdrawalWorker) runWorker(wChan chan *models.WithdrawalRecord, stop chan struct{}) {
	for {
		select {
		case wInfo := <-wChan:
			cww.do(wInfo)
		case <-stop:
			cww.wg.Done()
			return
		}
	}
}

func (cww *CashWithdrawalWorker) do(wInfo *models.WithdrawalRecord) {
	// get fx account info
	fxAccount := &models.FxAccount{
		ID: wInfo.AccountId,
	}
	has, err := models.GetFxAccountById(fxAccount)
	if err != nil {
		logrus.Errorf("get account[%d] info error: %v", wInfo.AccountId, err)
		return
	}
	if !has {
		logrus.Errorf("get account[%d] info error: no this account", wInfo.AccountId)
		return
	}

	if fxAccount.CanWithdrawals < wInfo.WithdrawalMoney {
		logrus.Errorf("withdrawal[%f] cannot be with in account[%v]", wInfo.WithdrawalMoney, fxAccount)
		return
	}

	err = models.MinusFxAccountMoney(wInfo.WithdrawalMoney, fxAccount)
	if err != nil {
		logrus.Errorf("withdrawal money[%f] with account[%v] error: %v", wInfo.WithdrawalMoney, fxAccount, err)
		return
	}

	// TODO: 给用户转账
	logrus.Infof("fx account[%v] withdrawal money[%f] update account success.", fxAccount, wInfo.WithdrawalMoney)

	wInfo.Status = WITHDRAWAL_DONE
	err = models.UpdateWithdrawalRecordStatus(wInfo)
	if err != nil {
		logrus.Errorf("withdrawal update record[%v] error: %v", wInfo, err)
		return
	}
}

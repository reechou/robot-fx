package act

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
	"github.com/reechou/robot-fx/logic/tools/order_check/fx_models"
)

type ActInfo struct {
	ActType        string
	ActHisType     int
	ActDesc        string
	ActValue       int64
	ActReward      float32
	ActUpperReward float32
	StartTime      int64
	EndTime        int64
}

type ActLogic struct {
	cfg *config.Config

	actList []*ActInfo
}

func NewActLogic(cfg *config.Config) *ActLogic {
	al := &ActLogic{
		cfg: cfg,
	}
	al.init()

	return al
}

func (self *ActLogic) init() {
	for _, v := range self.cfg.ActList {
		vlist := strings.Split(v, "|")
		if len(vlist) != 8 {
			logrus.Errorf("act[%s] error, %v.", v, vlist)
			continue
		}
		ai := &ActInfo{
			ActType: vlist[0],
			ActDesc: vlist[2],
		}
		t, err := strconv.Atoi(vlist[1])
		if err != nil {
			continue
		}
		ai.ActHisType = t
		value, err := strconv.ParseInt(vlist[3], 10, 0)
		if err != nil {
			continue
		}
		ai.ActValue = value
		reward, err := strconv.ParseFloat(vlist[4], 32)
		if err != nil {
			continue
		}
		ai.ActReward = float32(reward)
		supperReward, err := strconv.ParseFloat(vlist[5], 32)
		if err != nil {
			continue
		}
		ai.ActUpperReward = float32(supperReward)
		startTime, err := strconv.ParseInt(vlist[6], 10, 0)
		if err != nil {
			continue
		}
		ai.StartTime = startTime
		endTime, err := strconv.ParseInt(vlist[7], 10, 0)
		if err != nil {
			continue
		}
		ai.EndTime = endTime

		logrus.Debugf("load act: %v", ai)

		self.actList = append(self.actList, ai)
	}
}

func (self *ActLogic) CheckActOfOrder(fxAccount *fx_models.FxWxAccount, upperFxAccount *fx_models.FxWxAccount) {
	for _, v := range self.actList {
		count, err := self.checkOrderCount(fxAccount, v)
		if err != nil {
			logrus.Errorf("check account[%v] act[%v] order error: %v", fxAccount, v, err)
			continue
		}
		if v.ActValue == count {
			self.addActReward(fxAccount, upperFxAccount, v)
		}
	}
}

func (self *ActLogic) checkOrderCount(fxAccount *fx_models.FxWxAccount, act *ActInfo) (int64, error) {
	return fx_models.GetFxOrderSettlementRecordListCountById(fxAccount.ID, act.StartTime, act.EndTime)
}

func (self *ActLogic) addActReward(fxAccount *fx_models.FxWxAccount, upperFxAccount *fx_models.FxWxAccount, info *ActInfo) error {
	err := fx_models.AddFxWxAccountMoney(info.ActReward, fxAccount)
	if err != nil {
		logrus.Errorf("act[%s] add account[%v] money error: %v", info, fxAccount, err)
		return err
	}
	err = fx_models.AddFxWxAccountMoney(info.ActUpperReward, upperFxAccount)
	if err != nil {
		logrus.Errorf("act[%s] add upper account[%v] money error: %v", info, upperFxAccount, err)
		return err
	}
	self.addAccountHistory(fxAccount, upperFxAccount, info)
	return nil
}

func (self *ActLogic) addAccountHistory(fxAccount *fx_models.FxWxAccount, upperFxAccount *fx_models.FxWxAccount, info *ActInfo) {
	var historyList []fx_models.FxAccountHistory
	historyList = append(historyList, fx_models.FxAccountHistory{
		AccountId:  fxAccount.ID,
		UnionId:    fxAccount.WxId,
		Score:      info.ActReward,
		ChangeType: int64(info.ActHisType),
		ChangeDesc: info.ActDesc,
		CreatedAt:  time.Now().Unix(),
	})
	historyList = append(historyList, fx_models.FxAccountHistory{
		AccountId:  upperFxAccount.ID,
		UnionId:    upperFxAccount.WxId,
		Score:      info.ActUpperReward,
		ChangeType: int64(info.ActHisType),
		ChangeDesc: fmt.Sprintf("下线 %s %s", fxAccount.Name, info.ActDesc),
		CreatedAt:  time.Now().Unix(),
	})
	err := fx_models.CreateFxAccountHistoryList(historyList)
	if err != nil {
		logrus.Errorf("account[%v] upperAccount[%v] act[%v] fx account history list error: %v", fxAccount, upperFxAccount, info, err)
	}
}

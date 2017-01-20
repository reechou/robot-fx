package controller

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/models"
)

func (daemon *Daemon) CreateFxAccount(fxAccount *models.FxAccount) (int, error) {
	if fxAccount.Superior == "" {
		fxAccount.Superior = GodSalesman
	}

	if err := models.CreateFxAccount(fxAccount); err != nil {
		logrus.Errorf("create fx account error: %v", err)
		return 0, err
	}
	if fxAccount.Superior != "" && fxAccount.Superior != GodSalesman {
		superFxAccount := &models.FxAccount{
			WechatUnionId: fxAccount.Superior,
		}
		has, err := models.GetFxAccountFromWxUnionId(superFxAccount)
		if err != nil {
			logrus.Errorf("get super fx account error: %v", err)
			return 0, err
		}
		if !has {
			return 0, nil
		}
		err = models.AddFxAccountMoney(float32(daemon.cfg.Score.FollowScore), superFxAccount)
		if err != nil {
			logrus.Errorf("add super fx account money error: %v", err)
			return 0, err
		}
		h := models.FxAccountHistory{
			UnionId:    superFxAccount.UnionId,
			Score:      float32(daemon.cfg.Score.FollowScore),
			ChangeType: int64(FX_HISTORY_TYPE_INVITE),
			ChangeDesc: FxHistoryDescs[FX_HISTORY_TYPE_INVITE],
			CreatedAt:  time.Now().Unix(),
		}
		models.CreateFxAccountHistoryList([]models.FxAccountHistory{h})
		return daemon.cfg.Score.FollowScore, nil
	}

	return 0, nil
}

func (daemon *Daemon) CreateSalesman(fxAccount *models.FxAccount) error {
	return models.UpdateFxAccountSalesman(fxAccount)
}

func (daemon *Daemon) UpdateFxAccountBaseInfo(fxAccount *models.FxAccount) error {
	return models.UpdateFxAccountBaseInfo(fxAccount)
}

func (daemon *Daemon) UpdateFxAccountStatus(fxAccount *models.FxAccount) error {
	return models.UpdateFxAccountStatus(fxAccount)
}

func (daemon *Daemon) UpdateFxAccountSignTime(fxAccount *models.FxAccount) (int64, int, error) {
	affected, err := models.UpdateFxAccountSignTime(float32(daemon.cfg.Score.SignScore), fxAccount)
	if err != nil {
		return 0, 0, err
	}
	if affected > 0 {
		h := models.FxAccountHistory{
			UnionId:    fxAccount.UnionId,
			Score:      float32(daemon.cfg.Score.SignScore),
			ChangeType: int64(FX_HISTORY_TYPE_SIGN),
			ChangeDesc: FxHistoryDescs[FX_HISTORY_TYPE_SIGN],
			CreatedAt:  time.Now().Unix(),
		}
		models.CreateFxAccountHistoryList([]models.FxAccountHistory{h})
	}

	return affected, daemon.cfg.Score.SignScore, nil

}

func (daemon *Daemon) GetFxAccount(fxAccount *models.FxAccount) error {
	has, err := models.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("get fx account error: %v", err)
		return err
	}
	if !has {
		return fmt.Errorf("no this account.")
	}

	return nil
}

func (daemon *Daemon) GetFxAccountRank(offset, num int64) ([]models.FxAccount, error) {
	return models.GetFxAccountRank(offset, num)
}

func (daemon *Daemon) GetLowerPeopleCount(unionId string) (int64, error) {
	return models.GetLowerPeopleCount(unionId)
}

func (daemon *Daemon) GetLowerPeopleList(unionId string, offset, num int64) ([]models.FxAccount, error) {
	return models.GetLowerPeople(unionId, offset, num)
}

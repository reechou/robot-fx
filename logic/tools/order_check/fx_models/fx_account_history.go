package fx_models

import (
	"github.com/Sirupsen/logrus"
)

type FxAccountHistory struct {
	ID         int64   `xorm:"pk autoincr"`
	AccountId  int64   `xorm:"not null default 0 index"`
	UnionId    string  `xorm:"not null default '' varchar(128) index"`
	Score      float32 `xorm:"not null default 0.000 decimal(9,3)"`
	ChangeType int64   `xorm:"not null default 0 int index"`
	ChangeDesc string  `xorm:"not null default '' varchar(64)"`
	CreatedAt  int64   `xorm:"not null default 0 int"`
}

func CreateFxAccountHistoryList(list []FxAccountHistory) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		logrus.Errorf("create fx account history record list error: %v", err)
		return err
	}
	return nil
}

func GetFxAccountHistoryListCount(unionId string) (int64, error) {
	count, err := x.Where("union_id = ?", unionId).Count(&FxAccountHistory{})
	if err != nil {
		logrus.Errorf("union_id[%s] get fx account history list count error: %v", unionId, err)
		return 0, err
	}
	return count, nil
}

func GetFxAccountHistoryList(unionId string, offset, num int64) ([]FxAccountHistory, error) {
	var fxAccountHistoryList []FxAccountHistory
	err := x.Where("union_id = ?", unionId).Limit(int(num), int(offset)).Find(&fxAccountHistoryList)
	if err != nil {
		logrus.Errorf("union_id[%s] get fx account wait history list error: %v", unionId, err)
		return nil, err
	}
	return fxAccountHistoryList, nil
}

func GetFxAccountHistoryListByTypeCount(unionId string, cType int64) (int64, error) {
	count, err := x.Where("union_id = ?", unionId).And("change_type = ?", cType).Count(&FxAccountHistory{})
	if err != nil {
		logrus.Errorf("unionId[%s] get fx account history list by type[%d] count error: %v", unionId, cType, err)
		return 0, err
	}
	return count, nil
}

func GetFxAccountHistoryListByType(unionId string, cType, offset, num int64) ([]FxAccountHistory, error) {
	var fxAccountHistoryList []FxAccountHistory
	err := x.Where("union_id = ?", unionId).And("change_type = ?", cType).Limit(int(num), int(offset)).Find(&fxAccountHistoryList)
	if err != nil {
		logrus.Errorf("union_id[%s] get fx account wait history list by type[%d] error: %v", unionId, cType, err)
		return nil, err
	}
	return fxAccountHistoryList, nil
}

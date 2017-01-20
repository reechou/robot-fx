package controller

import (
	"github.com/reechou/robot-fx/logic/models"
)

func (daemon *Daemon) GetFxAccountHistoryListCount(unionId string) (int64, error) {
	return models.GetFxAccountHistoryListCount(unionId)
}

func (daemon *Daemon) GetFxAccountHistoryList(unionId string, offset, num int64) ([]models.FxAccountHistory, error) {
	return models.GetFxAccountHistoryList(unionId, offset, num)
}

func (daemon *Daemon) GetFxAccountHistoryListByTypeCount(unionId string, cType int64) (int64, error) {
	return models.GetFxAccountHistoryListByTypeCount(unionId, cType)
}

func (daemon *Daemon) GetFxAccountHistoryListByType(unionId string, cType, offset, num int64) ([]models.FxAccountHistory, error) {
	return models.GetFxAccountHistoryListByType(unionId, cType, offset, num)
}

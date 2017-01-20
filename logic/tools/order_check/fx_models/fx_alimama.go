package fx_models

import (
	"github.com/Sirupsen/logrus"
)

type FxRobotAlimama struct {
	ID        int64  `xorm:"pk autoincr"`
	RobotWX   string `xorm:"not null default '' varchar(128) index"`
	Alimama   string `xorm:"not null default '' varchar(128)"`
	CreatedAt int64  `xorm:"not null default 0 int"`
	UpdatedAt int64  `xorm:"not null default 0 int"`
}

func GetFxRobotAlimamaList() ([]FxRobotAlimama, error) {
	var fxRobotAlimamaList []FxRobotAlimama
	err := x.Where("id = 0").Find(&fxRobotAlimamaList)
	if err != nil {
		logrus.Errorf("get robot alimama list error: %v", err)
		return nil, err
	}
	return fxRobotAlimamaList, nil
}

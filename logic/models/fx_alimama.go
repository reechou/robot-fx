package models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type FxRobotAlimama struct {
	ID        int64  `xorm:"pk autoincr"`
	RobotWX   string `xorm:"not null default '' varchar(128) index"`
	Alimama   string `xorm:"not null default '' varchar(128)"`
	CreatedAt int64  `xorm:"not null default 0 int"`
	UpdatedAt int64  `xorm:"not null default 0 int"`
}

func CreateFxRobotAlimama(info *FxRobotAlimama) (err error) {
	if info.RobotWX == "" || info.Alimama == "" {
		return fmt.Errorf("fx robot[%s] alimama[%s] cannot be nil.", info.RobotWX, info.Alimama)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err = x.Insert(info)
	if err != nil {
		logrus.Errorf("create robot[%s] alimama[%s] error: %v", info.RobotWX, info.Alimama, err)
		return err
	}
	logrus.Infof("create robot[%s] alimama[%s] success.", info.RobotWX, info.Alimama)

	return
}

func GetRobotAlimama(info *FxRobotAlimama) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWX).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Errorf("cannot find alimama from robot[%s]", info.RobotWX)
		return false, nil
	}
	return true, nil
}

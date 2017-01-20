package models

import (
	"time"

	"github.com/Sirupsen/logrus"
)

type WithdrawalRecord struct {
	ID              int64   `xorm:"pk autoincr"`
	AccountId       int64   `xorm:"not null default 0 int index"`
	UnionId         string  `xorm:"not null default '' varchar(128) index"`
	RobotWx         string  `xorm:"not null default '' varchar(128)"` // 机器人昵称
	Wechat          string  `xorm:"not null default '' varchar(128)"` // 用户微信号
	Name            string  `xorm:"not null default '' varchar(128)"` // 用户昵称
	WithdrawalMoney float32 `xorm:"not null default 0.00 float(9,2)"`
	Balance         float32 `xorm:"not null default 0.00 float(9,2)"`
	OpenId          string  `xorm:"not null default '' varchar(128)"`
	Status          int64   `xorm:"not null default 0 int"`
	CreatedAt       int64   `xorm:"not null default 0 int"`
	UpdatedAt       int64   `xorm:"not null default 0 int index"`
}

type WithdrawalRecordError struct {
	ID              int64   `xorm:"pk autoincr"`
	AccountId       int64   `xorm:"not null default 0 int index"`
	UnionId         string  `xorm:"not null default '' varchar(128) index"`
	Name            string  `xorm:"not null default '' varchar(128) index"`
	WithdrawalMoney float32 `xorm:"not null default 0.00 float(9,2)"`
	ErrorMsg        string  `xorm:"not null default '' varchar(512)"`
	CreatedAt       int64   `xorm:"not null default 0 int index"`
}

func CreateWithdrawalRecord(info *WithdrawalRecord) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx withdrawal record union_id[%s] error: %v", info.UnionId, err)
		return err
	}
	logrus.Infof("create fx withdrawal record union_id[%s] withdrawal[%f] balance[%f] create success.", info.UnionId, info.WithdrawalMoney, info.Balance)
	return nil
}

func GetMonthWithdrawalRecord(unionId string) (int64, error) {
	timeStr := time.Now().Format("2006-01")
	t, _ := time.Parse("2006-01", timeStr)
	monthZero := t.Unix() - 8*3600
	count, err := x.Where("union_id = ?", unionId).And("updated_at > ?", monthZero).Count(&WithdrawalRecord{})
	if err != nil {
		logrus.Errorf("get month withdrawal record error: %v", err)
		return 0, err
	}
	return count, nil
}

func UpdateWithdrawalRecordStatus(info *WithdrawalRecord) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("status", "updated_at").Update(info, &WithdrawalRecord{ID: info.ID})
	return err
}

func GetWithdrawalRecordListCount(unionId string, status int64) (int64, error) {
	var count int64
	var err error
	if status == 0 {
		count, err = x.Where("union_id = ?", unionId).Count(&WithdrawalRecord{})
	} else {
		count, err = x.Where("union_id = ?", unionId).And("status = ?", status).Count(&WithdrawalRecord{})
	}
	if err != nil {
		logrus.Errorf("union_id[%s] get withdrawal record list count error: %v", unionId, err)
		return 0, err
	}
	return count, nil
}

func GetWithdrawalRecordListCountById(accountId int64) (int64, error) {
	count, err := x.Where("account_id = ?", accountId).Count(&WithdrawalRecord{})
	if err != nil {
		logrus.Errorf("account_id[%d] get withdrawal record list count error: %v", accountId, err)
		return 0, err
	}
	return count, nil
}

func GetWithdrawalRecordList(unionId string, offset, num, status int64) ([]WithdrawalRecord, error) {
	var list []WithdrawalRecord
	var err error
	if status == 0 {
		err = x.Where("union_id = ?", unionId).Limit(int(num), int(offset)).Find(&list)
	} else {
		err = x.Where("union_id = ?", unionId).And("status = ?", status).Limit(int(num), int(offset)).Find(&list)
	}
	if err != nil {
		logrus.Errorf("union_id[%s] get withdrawal record list error: %v", unionId, err)
		return nil, err
	}
	return list, nil
}

func GetWithdrawalRecordListById(accountId int64, offset, num int64) ([]WithdrawalRecord, error) {
	var list []WithdrawalRecord
	err := x.Where("account_id = ?", accountId).Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		logrus.Errorf("account_id[%d] get withdrawal record list error: %v", accountId, err)
		return nil, err
	}
	return list, nil
}

func GetWithdrawalRecordSum(unionId string) (float32, error) {
	total, err := x.Where("union_id = ?", unionId).Sum(&WithdrawalRecord{}, "withdrawal_money")
	if err != nil {
		logrus.Errorf("unionId[%s] get withdrawal sum error: %v", unionId, err)
		return 0.0, err
	}
	return float32(total), nil
}

func CreateWithdrawalRecordError(info *WithdrawalRecordError) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx withdrawal error msg record[%v] error: %v", info, err)
		return err
	}
	logrus.Infof("create fx withdrawal error msg record[%v] create success.", info)
	return nil
}

func GetWithdrawalRecordErrorListCount() (int64, error) {
	count, err := x.Count(&WithdrawalRecordError{})
	if err != nil {
		logrus.Errorf("get withdrawal error msg record list count error: %v", err)
		return 0, err
	}
	return count, nil
}

func GetWithdrawalRecordErrorList(offset, num int64) ([]WithdrawalRecordError, error) {
	var list []WithdrawalRecordError
	var err error
	err = x.Desc("created_at").Limit(int(num), int(offset)).Find(&list)
	if err != nil {
		logrus.Errorf("get withdrawal error msg record list error: %v", err)
		return nil, err
	}
	return list, nil
}

func GetWithdrawalRecordErrorListFromName(name string) ([]WithdrawalRecordError, error) {
	var list []WithdrawalRecordError
	var err error
	err = x.Where("name like ?", name).Find(&list)
	//results, err := x.Query("select * from withdrawal_record_error where name like '%?%'", name)
	if err != nil {
		logrus.Errorf("get withdrawal error msg record list from name[%s] error: %v", name, err)
		return nil, err
	}
	return list, nil
}

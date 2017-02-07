package models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type FxAccount struct {
	ID             int64   `xorm:"pk autoincr"`
	UnionId        string  `xorm:"not null default '' varchar(128) unique"` // robotwx_$$_nickname - 可变
	RobotWx        string  `xorm:"not null default '' varchar(128) index"`  // 机器人微信昵称 - 可控不可变
	WechatUnionId  string  `xorm:"not null default '' varchar(128) unique"` // robotwx_$$_wxid - 不可变
	Wechat         string  `xorm:"not null default '' varchar(128)"`        // 用户微信号,如果没有则为 wxid
	WxId           string  `xorm:"not null default '' varchar(128) index"`  // 用户wxid
	UserName       string  `xorm:"not null default '' varchar(128) index"`  // robot本次登录状态下唯一id
	Name           string  `xorm:"not null default '' varchar(128)"`        // 用户昵称 - 可变
	CanWithdrawals float32 `xorm:"not null default 0.00 float(9,2)"`        // 可提现
	AllScore       float32 `xorm:"not null default 0.00 float(9,2)"`        // 总积分
	Superior       string  `xorm:"not null default '' varchar(128) index"`  // 上级
	MemberId       string  `xorm:"not null default '' varchar(32)"`         // 阿里妈妈用户id
	GuideId        string  `xorm:"not null default '' varchar(32)"`         // 导购位
	AdzoneId       string  `xorm:"not null default '' varchar(32)"`         // 广告位
	SignTime       int64   `xorm:"not null default 0 int index"`
	Status         int64   `xorm:"not null default 0 int"`
	CreatedAt      int64   `xorm:"not null default 0 int index"`
	UpdatedAt      int64   `xorm:"not null default 0 int"`
}

type FxWxAccount struct {
	ID             int64   `xorm:"pk autoincr"`
	Wechat         string  `xorm:"not null default '' varchar(128)"`        // 用户微信号,如果没有则为 wxid
	WxId           string  `xorm:"not null default '' varchar(128) unique"` // 用户wxid
	Name           string  `xorm:"not null default '' varchar(128)"`        // 用户昵称 - 可变
	CanWithdrawals float32 `xorm:"not null default 0.00 float(9,2)"`        // 可提现
	AllScore       float32 `xorm:"not null default 0.00 float(9,2)"`        // 总积分
	Superior       string  `xorm:"not null default '' varchar(128) index"`  // 上级
	SignTime       int64   `xorm:"not null default 0 int index"`
	CreatedAt      int64   `xorm:"not null default 0 int index"`
	UpdatedAt      int64   `xorm:"not null default 0 int"`
}

func CreateFxWxAccount(info *FxWxAccount) error {
	if info.WxId == "" {
		return fmt.Errorf("wx wxid[%s] cannot be nil.", info.WxId)
	}
	
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx wx account error: %v", err)
		return err
	}
	logrus.Infof("create fx wx account from wxid[%s] success.", info.WxId)
	
	return nil
}

func GetFxWxAccount(info *FxWxAccount) (bool, error) {
	has, err := x.Where("wx_id = ?", info.WxId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Debugf("cannot find fx wx account from wxid[%s]", info.WxId)
		return false, nil
	}
	return true, nil
}

func AddFxWxAccountMoney(allAdd float32, info *FxWxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	_, err = x.Exec("update fx_wx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=? where wx_id=?",
		allAdd, allAdd, info.UpdatedAt, info.WxId)
	if err != nil {
		return err
	}
	logrus.Infof("fx wx account[%s] add money[%f] success.", info.WxId, allAdd)
	return nil
}

func MinusFxWxAccountMoney(allMinus float32, info *FxWxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	result, err := x.Exec("update fx_wx_account set can_withdrawals=can_withdrawals-?, updated_at=? where wx_id=? and can_withdrawals >= ?",
		allMinus, info.UpdatedAt, info.WxId, allMinus)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		logrus.Errorf("minus fx wx account error affected == 0")
		return fmt.Errorf("minus fx wx account error affected == 0")
	}
	logrus.Infof("fx wx account[%s] minus money[%f] success.", info.WxId, allMinus)
	return nil
}

func UpdateFxWxAccountSignTime(allAdd float32, info *FxWxAccount) (int64, error) {
	now := time.Now().Unix()
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	dayZero := t.Unix() - 8*3600
	result, err := x.Exec("update fx_wx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=?, sign_time=? where wx_id=? and sign_time < ?",
		allAdd, allAdd, now, now, info.WxId, dayZero)
	if err != nil {
		logrus.Errorf("update fx_account sign time error: %v", err)
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("get affected error: %v", err)
		return 0, err
	}
	return affected, nil
}

func GetFxWxLowerPeopleCount(wxid string) (int64, error) {
	count, err := x.Where("superior = ?", wxid).Count(&FxWxAccount{})
	if err != nil {
		logrus.Errorf("wxid[%s] get lower peoples list count error: %v", wxid, err)
		return 0, err
	}
	return count, nil
}

func GetFxWxLowerPeople(wxid string, offset, num int64) ([]FxWxAccount, error) {
	var lowerPeoples []FxWxAccount
	err := x.Where("superior = ?", wxid).Desc("created_at").Limit(int(num), int(offset)).Find(&lowerPeoples)
	if err != nil {
		logrus.Errorf("wxid[%s] lower peoples list error: %v", wxid, err)
		return nil, err
	}
	return lowerPeoples, nil
}


func CreateFxAccount(info *FxAccount) (err error) {
	if info.UnionId == "" {
		return fmt.Errorf("wx union_id[%s] cannot be nil.", info.UnionId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err = x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx account error: %v", err)
		return err
	}
	logrus.Infof("create fx account from wx_unionid[%s] success.", info.UnionId)

	return
}

func UpdateFxAccountBaseInfo(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("phone", "name", "updated_at").Update(info, &FxAccount{UnionId: info.UnionId})
	return err
}

func UpdateFxAccountStatus(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("status", "updated_at").Update(info, &FxAccount{UnionId: info.UnionId})
	return err
}

func UpdateFxAccountSalesman(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("ticket", "phone", "updated_at").Update(info, &FxAccount{UnionId: info.UnionId})
	return err
}

func UpdateFxAccountName(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("union_id", "name", "updated_at").Update(info, &FxAccount{RobotWx: info.RobotWx, UserName: info.UserName})
	return err
}

func UpdateFxAccountUserName(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("user_name", "updated_at").Update(info, &FxAccount{UnionId: info.UnionId})
	return err
}

func UpdateFxAccountWechat(info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Cols("wechat_union_id", "member_id", "guide_id", "adzone_id", "updated_at").Update(info, &FxAccount{UnionId: info.UnionId})
	return err
}

func UpdateFxAccountSignTime(allAdd float32, info *FxAccount) (int64, error) {
	now := time.Now().Unix()
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.Parse("2006-01-02", timeStr)
	dayZero := t.Unix() - 8*3600
	result, err := x.Exec("update fx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=?, sign_time=? where union_id=? and sign_time < ?",
		allAdd, allAdd, now, now, info.UnionId, dayZero)
	if err != nil {
		logrus.Errorf("update fx_account sign time error: %v", err)
		return 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("get affected error: %v", err)
		return 0, err
	}
	return affected, nil
}

func AddFxAccountMoney(allAdd float32, info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	_, err = x.Exec("update fx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=? where union_id=?",
		allAdd, allAdd, info.UpdatedAt, info.UnionId)
	if err != nil {
		return err
	}
	logrus.Infof("fx account[%s] add money[%f] success.", info.UnionId, allAdd)
	return nil
}

func AddFxAccountMoneyWithWxUnionId(allAdd float32, info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	_, err = x.Exec("update fx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=? where wechat_union_id=?",
		allAdd, allAdd, info.UpdatedAt, info.WechatUnionId)
	if err != nil {
		return err
	}
	logrus.Infof("fx account[%s] add money[%f] success.", info.UnionId, allAdd)
	return nil
}

func MinusFxAccountMoney(allMinus float32, info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	_, err = x.Exec("update fx_account set can_withdrawals=can_withdrawals-?, updated_at=? where union_id=? and can_withdrawals >= ?",
		allMinus, info.UpdatedAt, info.UnionId, allMinus)
	if err != nil {
		return err
	}
	logrus.Infof("fx account[%s] minus money[%f] success.", info.UnionId, allMinus)
	return nil
}

func GetFxAccount(info *FxAccount) (bool, error) {
	has, err := x.Where("union_id = ?", info.UnionId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Debugf("cannot find fx account from unionid[%s]", info.UnionId)
		return false, nil
	}
	return true, nil
}

func GetFxAccountFromUserName(info *FxAccount) (bool, error) {
	has, err := x.Where("robot_wx = ?", info.RobotWx).And("user_name = ?", info.UserName).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Debugf("cannot find fx account from user_name[%s]", info.UserName)
		return false, nil
	}
	return true, nil
}

func GetFxAccountFromWxUnionId(info *FxAccount) (bool, error) {
	has, err := x.Where("wechat_union_id = ?", info.WechatUnionId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Debugf("cannot find fx account from wechat_union_id[%s]", info.WechatUnionId)
		return false, nil
	}
	return true, nil
}

func GetFxAccountById(info *FxAccount) (bool, error) {
	has, err := x.Where("id = ?", info.UnionId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Errorf("cannot find fx account from account_id[%s]", info.ID)
		return false, nil
	}
	return true, nil
}

func GetLowerPeopleCount(unionId string) (int64, error) {
	count, err := x.Where("superior = ?", unionId).Count(&FxAccount{})
	if err != nil {
		logrus.Errorf("union_id[%s] get lower peoples list count error: %v", unionId, err)
		return 0, err
	}
	return count, nil
}

func GetFxAccountRank(offset, num int64) ([]FxAccount, error) {
	var rankList []FxAccount
	err := x.Where("id != 12").Desc("all_score").Limit(int(num), int(offset)).Find(&rankList)
	if err != nil {
		logrus.Errorf("get fx account rank list error: %v", err)
		return nil, err
	}
	return rankList, nil
}

func GetLowerPeople(unionId string, offset, num int64) ([]FxAccount, error) {
	var lowerPeoples []FxAccount
	err := x.Where("superior = ?", unionId).Desc("created_at").Limit(int(num), int(offset)).Find(&lowerPeoples)
	if err != nil {
		logrus.Errorf("union_id[%s] lower peoples list error: %v", unionId, err)
		return nil, err
	}
	return lowerPeoples, nil
}

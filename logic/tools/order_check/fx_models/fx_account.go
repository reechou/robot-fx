package fx_models

import (
	"time"

	"github.com/Sirupsen/logrus"
)

type FxAccount struct {
	ID             int64   `xorm:"pk autoincr"`
	UnionId        string  `xorm:"not null default '' varchar(128) unique"` // robotwx_$$_nickname - 可变
	RobotWx        string  `xorm:"not null default '' varchar(128) index"`  // 机器人微信昵称 - 可控不可变
	WechatUnionId  string  `xorm:"not null default '' varchar(128) unique"` // robotwx_$$_wechat - 不可变
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

func AddFxAccountMoney(allAdd float32, info *FxAccount) error {
	info.UpdatedAt = time.Now().Unix()
	var err error
	_, err = x.Exec("update fx_account set can_withdrawals=can_withdrawals+?, all_score=all_score+?, updated_at=? where union_id=?",
		allAdd, allAdd, info.UpdatedAt, info.UnionId)
	return err
}

func GetFxAccount(info *FxAccount) (bool, error) {
	has, err := x.Where("union_id = ?", info.UnionId).Get(info)
	if err != nil {
		return false, err
	}
	if !has {
		logrus.Debugf("cannot find fx account from wx_unionid[%s]", info.UnionId)
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

package fx_models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
)

type FxOrder struct {
	ID          int64   `xorm:"pk autoincr"`
	AccountId   int64   `xorm:"not null default 0 int index"`
	UnionId     string  `xorm:"not null default '' varchar(128) index"`
	OrderId     string  `xorm:"varchar(128) not null default '' unique(uni_fx_order_id)"`
	GoodsId     string  `xorm:"varchar(128) not null default '' unique(uni_fx_order_id)"`
	OrderName   string  `xorm:"not null default '' varchar(128)"`
	Price       float32 `xorm:"not null default 0.00 float(9,2) unique(uni_fx_order_id)"`
	ReturnMoney float32 `xorm:"not null default 0.00 float(9,2)" json:"-"`
	AdName      string  `xorm:"varchar(128) not null default ''"`
	Status      int64   `xorm:"not null default 0 int index"`
	CreatedAt   int64   `xorm:"not null default 0 int index"`
	UpdatedAt   int64   `xorm:"not null default 0 int"`
}

type FxOrderWaitSettlementRecord struct {
	ID          int64   `xorm:"pk autoincr"`
	AccountId   int64   `xorm:"not null default 0 int index"`
	UnionId     string  `xorm:"not null default '' varchar(128) index"`
	OrderId     string  `xorm:"not null default '' varchar(128) index"`
	GoodsId     string  `xorm:"not null default '' varchar(128)"`
	Price       float32 `xorm:"not null default 0.00 float(9,2)"`
	ReturnMoney float32 `xorm:"not null default 0.00 float(9,2)"`
	Level       int64   `xorm:"not null default 0 int index"`
	CreatedAt   int64   `xorm:"not null default 0 int index"`
}

type FxOrderSettlementRecord struct {
	ID           int64   `xorm:"pk autoincr"`
	AccountId    int64   `xorm:"not null default 0 int index"`
	UnionId      string  `xorm:"not null default '' varchar(128) index"`
	OrderId      string  `xorm:"not null default '' varchar(128)"`
	GoodsId      string  `xorm:"not null default '' varchar(128)"`
	Price        float32 `xorm:"not null default 0.00 float(9,2)"`
	ReturnMoney  float32 `xorm:"not null default 0.00 float(9,2)"`
	SourceId     string  `xorm:"not null default '' varchar(128)"`
	Level        int64   `xorm:"not null default 0 int index"`
	OrderCreated int64   `xorm:"not null default 0 int index"`
	CreatedAt    int64   `xorm:"not null default 0 int index"`
	UpdatedAt    int64   `xorm:"not null default 0 int"`
}

func CreateFxOrder(info *FxOrder) error {
	if info.UnionId == "" || info.OrderId == "" || info.GoodsId == "" {
		return fmt.Errorf("fx order union_id[%s] order_id[%s] goods_id[%s] cannot be nil.", info.UnionId, info.OrderId, info.GoodsId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx order error: %v", err)
		return err
	}
	logrus.Infof("fx order union_id[%s] order_id[%s] create success.", info.UnionId, info.OrderId)

	return nil
}

func GetFxOrderInfo(info *FxOrder) (bool, error) {
	has, err := x.Where("order_id = ?", info.OrderId).And("goods_id = ?", info.GoodsId).And("price = ?", info.Price).Get(info)
	if err != nil {
		logrus.Errorf("get fx order[%s] error: %v", info.OrderId, err)
		return false, fmt.Errorf("get fx order[%s] error: %v", info.OrderId, err)
	}
	if !has {
		logrus.Debugf("get fx order[%s] has no this order.", info.OrderId)
		return false, nil
	}
	return true, nil
}

func UpdateFxOrderStatus(info *FxOrder) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Id(info.ID).Cols("status", "updated_at").Update(info)
	return err
}

func CreateFxOrderWaitSettlementRecordList(list []FxOrderWaitSettlementRecord) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		logrus.Errorf("create fx order wait settlement record list error: %v", err)
		return err
	}
	return nil
}

func IterateFxWaitOrder(status int64, f xorm.IterFunc) error {
	logrus.Debugf("interate fx waiting orders...")
	err := x.Where("status = ?", status).Iterate(&FxOrder{}, f)
	if err != nil {
		logrus.Errorf("iterate fx order error: %v", err)
		return err
	}
	return nil
}

func CreateFxOrderSettlementRecordList(list []FxOrderSettlementRecord) error {
	if len(list) == 0 {
		return nil
	}
	_, err := x.Insert(&list)
	if err != nil {
		logrus.Errorf("create fx order settlement record list error: %v", err)
		return err
	}
	return nil
}

func GetFxOrderSettlementRecordListCountById(accountId, startTime, endTime int64) (int64, error) {
	count, err := x.Where("account_id = ?", accountId).And("level = 0").And("order_created >= ?", startTime).And("order_created <= ?", endTime).Count(&FxOrderSettlementRecord{})
	if err != nil {
		logrus.Errorf("account_id[%d] get fx order settlement record list count error: %v", accountId, err)
		return 0, err
	}
	return count, nil
}

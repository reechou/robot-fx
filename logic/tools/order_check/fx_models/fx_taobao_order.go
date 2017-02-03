package fx_models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type TaobaoOrder struct {
	ID                  int64   `xorm:"pk autoincr"`
	OrderCreatedTime    string  `xorm:"not null default '' varchar(32)"` // 创建时间
	OrderClickTime      string  `xorm:"not null default '' varchar(32)"` // 点击时间
	GoodsInfo           string  `xorm:"not null default '' varchar(128)"`
	GoodsId             string  `xorm:"not null default '' varchar(64) unique(uni_taobao_order_id)"`
	WangwangName        string  `xorm:"not null default '' varchar(128)"`
	StoreName           string  `xorm:"not null default '' varchar(128)"`
	GoodsNum            int     `xorm:"not null default 0 int"`
	GoodsPrice          float32 `xorm:"not null default 0.00 float(9,2)"`
	OrderState          int     `xorm:"not null default 0 int"`
	OrderType           string  `xorm:"not null default '' varchar(32)"`
	IncomeRatio         float32 `xorm:"not null default 0.00 float(9,2)"`                             // 收入比率
	SplitRatio          float32 `xorm:"not null default 0.00 float(9,2)"`                             // 分成比率
	PayPrice            float32 `xorm:"not null default 0.00 float(9,2) unique(uni_taobao_order_id)"` // 付款金额
	PredictingEffect    float32 `xorm:"not null default 0.00 float(9,2)"`                             // 预估效果
	SettlementMoney     float32 `xorm:"not null default 0.00 float(9,2)"`                             // 结算金额
	EstimatedIncome     float32 `xorm:"not null default 0.00 float(9,2)"`                             // 预估收入
	SettlementTime      string  `xorm:"not null default '' varchar(32)"`
	CommissionRate      float32 `xorm:"not null default 0.00 float(9,2)"`
	CommissionMoney     float32 `xorm:"not null default 0.00 float(9,2)"`
	SubsidyRate         float32 `xorm:"not null default 0.00 float(9,2)"`
	SubsidyMoney        float32 `xorm:"not null default 0.00 float(9,2)"`
	SubsidyType         string  `xorm:"not null default '' varchar(64)"`
	TransactionPlatform string  `xorm:"not null default '' varchar(64)"`
	ThirdPartyService   string  `xorm:"not null default '' varchar(128)"`
	OrderId             string  `xorm:"not null default '' varchar(64) unique(uni_taobao_order_id)"`
	CategoryName        string  `xorm:"not null default '' varchar(128)"`
	SiteId              string  `xorm:"not null default '' varchar(16)"`
	SiteName            string  `xorm:"not null default '' varchar(128)"`
	AdId                string  `xorm:"not null default '' varchar(16)"`
	AdName              string  `xorm:"not null default '' varchar(128)"`
	CreatedAt           int64   `xorm:"not null default 0 int"`
	UpdatedAt           int64   `xorm:"not null default 0 int"`
}

func CreateTaobaoOrder(info *TaobaoOrder) error {
	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now

	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create taobao order error: %v", err)
		return err
	}
	logrus.Infof("create taobao order id[%s] success.", info.OrderId)

	return nil
}

func GetTaobaoOrder(info *TaobaoOrder) (bool, error) {
	has, err := x.Where("goods_id = ?", info.GoodsId).And("pay_price = ?", info.PayPrice).And("order_id = ?", info.OrderId).Get(info)
	if err != nil {
		logrus.Errorf("get taobao order[%s] error: %v", info.OrderId, err)
		return false, fmt.Errorf("get taobao order[%s] error: %v", info.OrderId, err)
	}
	if !has {
		logrus.Errorf("get taobao order[%s] has no this order.", info.OrderId)
		return false, nil
	}
	return true, nil
}

func UpdateTaobaoOrderStatus(info *TaobaoOrder) error {
	info.UpdatedAt = time.Now().Unix()
	_, err := x.Id(info.ID).Cols("order_state", "updated_at").Update(info)
	return err
}

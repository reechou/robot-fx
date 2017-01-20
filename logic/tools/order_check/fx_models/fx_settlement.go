package fx_models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type SettlementFxOrderInfo struct {
	Status        int64
	Order         *FxOrder
	OrderAddMoney float32
}

func SettlementOwnerFxOrder(info *SettlementFxOrderInfo) error {
	now := time.Now().Unix()
	result, err := x.Exec("update fx_account fa, fx_order fo set fa.can_withdrawals=fa.can_withdrawals+?, fa.all_score=fa.all_score+?, fa.updated_at=?, fo.status=?, fo.updated_at=? where fa.union_id=? and fo.order_id=? and fo.goods_id = ? and fo.price = ?",
		info.OrderAddMoney, info.OrderAddMoney, now, info.Status, now, info.Order.UnionId, info.Order.OrderId, info.Order.GoodsId, info.Order.Price)
	if err != nil {
		logrus.Errorf("settlement owner[%s] fx order[%s] error: %v", info.Order.UnionId, info.Order.OrderId, err)
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		logrus.Errorf("settlement get affected error: %v", err)
		return err
	}
	if affected == 0 {
		logrus.Errorf("settlement update error affected == 0")
		return fmt.Errorf("settlement update error affected == 0")
	}
	return nil
}

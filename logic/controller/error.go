package controller

import (
	"errors"
)

var (
	ErrWithdrawalOverMonthLimit = errors.New("超过当月提现限制")
	ErrWithdrawalMinimum        = errors.New("提现未到最低额度")
	ErrWithdrawalLimitBalance   = errors.New("提现金额超过余额")
)

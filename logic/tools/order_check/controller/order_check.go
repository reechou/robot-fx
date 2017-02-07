package controller

import (
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
	"github.com/reechou/robot-fx/logic/tools/order_check/fx_models"
	"github.com/reechou/robot-fx/utils"
)

type OrderCheck struct {
	cfg *config.Config

	sw  *SettlementWorker
	toc *TaobaoOrderCheck
	fom *FxOrderManager
	ohs *OrderHttpSrv

	wg   sync.WaitGroup
	stop chan struct{}
	done chan struct{}
}

func NewOrderCheck(cfg *config.Config) *OrderCheck {
	if cfg.Debug {
		utils.EnableDebug()
	}

	ocw := &OrderCheck{
		cfg:  cfg,
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	ocw.sw = NewSettlementWorker(cfg.WorkerInfo.SWMaxWorker, cfg.WorkerInfo.SWMaxChanLen, cfg)
	ocw.fom = NewFxOrderManager(cfg)
	ocw.toc = NewTaobaoOrderCheck(cfg, ocw.fom)
	ocw.ohs = NewOrderHTTPServer(cfg, ocw.toc)

	fx_models.InitDB(cfg)

	return ocw
}

func (ocw *OrderCheck) Stop() {
	close(ocw.stop)
	<-ocw.done
}

func (ocw *OrderCheck) Run() {
	go ocw.runLoop()
	
	ocw.ohs.Run()
}

func (ocw *OrderCheck) runLoop() {
	logrus.Debugf("start run fx order check...")
	
	ocw.runCheck()
	for {
		select {
		case <-time.After(time.Duration(ocw.cfg.WorkerInfo.OrderCheckInterval) * time.Second):
			ocw.runCheck()
		case <-ocw.stop:
			close(ocw.done)
			return
		}
	}
}

func (ocw *OrderCheck) runCheck() {
	err := fx_models.IterateFxWaitOrder(FX_ORDER_WAIT, ocw.handleOrder)
	if err != nil {
		logrus.Errorf("run check error: %v", err)
	}
}

func (ocw *OrderCheck) handleOrder(idx int, bean interface{}) error {
	order := bean.(*fx_models.FxOrder)
	logrus.Debugf("order[%v] checking.", order)
	// check order status
	taobaoOrder := &fx_models.TaobaoOrder{
		OrderId:  order.OrderId,
		GoodsId:  order.GoodsId,
		PayPrice: order.Price,
	}
	has, err := fx_models.GetTaobaoOrder(taobaoOrder)
	if err != nil {
		logrus.Errorf("get taobao order[%v] error: %v", order, err)
		return err
	}
	if !has {
		logrus.Errorf("get taobao order no this order[%v]", order)
		return nil
	}

	logrus.Debugf("get taobao order[%v]", taobaoOrder)

	if taobaoOrder.OrderState == TAOBAO_ORDER_SETTLEMENT || taobaoOrder.OrderState == TAOBAO_ORDER_SUCCESS {
		// do settlement
		ocw.sw.SettlementOrder(order)
	} else if taobaoOrder.OrderState == TAOBAO_ORDER_INVALID {
		order.Status = FX_ORDER_FAILED
		err := fx_models.UpdateFxOrderStatus(order)
		if err != nil {
			logrus.Errorf("order[%s] update status error.", order.OrderId)
		}
	}

	return nil
}

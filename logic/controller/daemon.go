package controller

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/reechou/robot-fx/config"
	"github.com/reechou/robot-fx/logic/ext"
)

type Daemon struct {
	cfg *config.Config

	client *http.Client
	r      *rand.Rand
	we     *ext.WechatExt
	dme    *ext.DuobbManagerExt
	cww    *CashWithdrawalWorker
}

func NewDaemon(cfg *config.Config) *Daemon {
	d := &Daemon{
		cfg:    cfg,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
		client: &http.Client{},
	}
	d.cww = NewCashWithdrawalWorker(DEFAULT_MAX_WORKER, DEFAULT_MAX_CHAN_LEN, d.cfg)
	d.we = ext.NewWechatExt(d.cfg)
	d.dme = ext.NewDuobbManagerExt(d.cfg)
	return d
}

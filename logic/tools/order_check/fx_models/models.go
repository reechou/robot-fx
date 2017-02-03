package fx_models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
)

var x *xorm.Engine

func InitDB(cfg *config.Config) {
	var err error
	x, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		cfg.FxDBInfo.User,
		cfg.FxDBInfo.Pass,
		cfg.FxDBInfo.Host,
		cfg.FxDBInfo.DBName))
	if err != nil {
		logrus.Fatalf("Fail to init db engine: %v", err)
	}
	x.SetLogger(nil)
	x.SetMapper(core.GonicMapper{})
	x.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	if err = x.Sync2(new(TaobaoOrder)); err != nil {
		logrus.Fatalf("Fail to sync database: %v", err)
	}
}

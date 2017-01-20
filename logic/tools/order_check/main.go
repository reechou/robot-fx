package main

import (
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
	"github.com/reechou/robot-fx/logic/tools/order_check/controller"
)

func main() {
	controller.NewOrderCheck(config.NewConfig()).Run()
}

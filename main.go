package main

import (
	"github.com/reechou/robot-fx/config"
	"github.com/reechou/robot-fx/servermain"
)

func main() {
	servermain.NewMain(config.NewConfig()).Run()
}

package logic

import (
	"time"

	"github.com/reechou/robot-fx/config"
	"github.com/reechou/robot-fx/server"
)

type LogicThinking interface {
	InitRouter(s *server.Server)
	Shutdown(timeout time.Duration)
}

func NewLogic(cfg *config.Config) (LogicThinking, error) {
	l := NewModuleLogic(cfg)
	return l, nil
}

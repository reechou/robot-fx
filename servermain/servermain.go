package servermain

import (
	"runtime"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/config"
	"github.com/reechou/robot-fx/logic"
	"github.com/reechou/robot-fx/pkg/listener"
	"github.com/reechou/robot-fx/pkg/pidfile"
	"github.com/reechou/robot-fx/pkg/signal"
	"github.com/reechou/robot-fx/server"
	"github.com/reechou/robot-fx/utils"
)

type ServerMain struct {
	cfg *config.Config
}

func NewMain(cfg *config.Config) *ServerMain {
	return &ServerMain{
		cfg: cfg,
	}
}

func (srv *ServerMain) Run() {
	GoMaxProcs := runtime.GOMAXPROCS(0)
	logrus.Infof("setting maximum number of CPUs to %d, total number of available CPUs is %d", GoMaxProcs, runtime.NumCPU())

	if srv.cfg.Debug {
		utils.EnableDebug()
	}

	var pfile *pidfile.PIDFile
	if srv.cfg.PidFile != "" {
		pf, err := pidfile.New(srv.cfg.PidFile)
		if err != nil {
			logrus.Fatalf("Error starting daemon: %v", err)
		}
		pfile = pf
		defer func() {
			if err := pfile.Remove(); err != nil {
				logrus.Error(err)
			}
		}()
	}

	if len(srv.cfg.Hosts) == 0 {
		srv.cfg.Hosts = make([]string, 1)
	}

	s := server.New(srv.cfg)

	for i := 0; i < len(srv.cfg.Hosts); i++ {
		protoAddr := srv.cfg.Hosts[i]
		protoAddrParts := strings.SplitN(protoAddr, "://", 2)
		if len(protoAddrParts) != 2 {
			logrus.Fatalf("bad format %s, expected PROTO://ADDR", protoAddr)
		}
		l, err := listener.Init(protoAddrParts[1], srv.cfg.TLSConfig)
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Debugf("Listener created for HTTP on %s (%s)", protoAddrParts[0], protoAddrParts[1])
		s.Accept(protoAddrParts[0], l...)
	}

	l, err := logic.NewLogic(srv.cfg)
	if err != nil || l == nil {
		if pfile != nil {
			if err := pfile.Remove(); err != nil {
				logrus.Error(err)
			}
		}
		logrus.Fatalf("Error starting server: %v", err)
	}
	logrus.Info("Server logic has completed initialization")

	l.InitRouter(s)

	serveWait := make(chan error)
	go s.Wait(serveWait)

	signal.Trap(func() {
		s.Close()
		<-serveWait
		// logic shutdown here
		l.Shutdown(15)
		if pfile != nil {
			if err := pfile.Remove(); err != nil {
				logrus.Error(err)
			}
		}
	})

	errServer := <-serveWait
	if errServer != nil {
		if pfile != nil {
			if err := pfile.Remove(); err != nil {
				logrus.Error(err)
			}
		}
		logrus.Fatalf("Shutting down due to serve error: %v", errServer)
	}
	logrus.Infof("Shutdown server success.")
}

package controller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
)

type OrderHttpSrv struct {
	cfg          *config.Config
	httpSrv      *HttpSrv
	orderManager *TaobaoOrderCheck
}

type HttpHandler func(rsp http.ResponseWriter, req *http.Request) (interface{}, error)

func NewOrderHTTPServer(cfg *config.Config, orderManager *TaobaoOrderCheck) *OrderHttpSrv {
	srv := &OrderHttpSrv{
		cfg:          cfg,
		orderManager: orderManager,
		httpSrv:      &HttpSrv{Host: cfg.OrderServerHost, Routers: make(map[string]http.HandlerFunc)},
	}
	srv.registerHandlers()

	return srv
}

func (self *OrderHttpSrv) Run() {
	logrus.Infof("order http server starting...")
	self.httpSrv.Run()
}

func (self *OrderHttpSrv) registerHandlers() {
	self.httpSrv.Route("/", self.Index)
	
	self.httpSrv.Route("/taobaoorder", self.httpWrap(self.TaobaoOrder))
}

func (self *OrderHttpSrv) httpWrap(handler HttpHandler) func(rsp http.ResponseWriter, req *http.Request) {
	f := func(rsp http.ResponseWriter, req *http.Request) {
		logURL := req.URL.String()
		start := time.Now()
		defer func() {
			logrus.Debugf("[OrderHttpSrv][httpWrap] http: request url[%s] use_time[%v]", logURL, time.Now().Sub(start))
		}()
		obj, err := handler(rsp, req)
		// check err
	HAS_ERR:
		if err != nil {
			logrus.Errorf("[OrderHttpSrv][httpWrap] http: request url[%s] error: %v", logURL, err)
			code := 500
			errMsg := err.Error()
			if strings.Contains(errMsg, "Permission denied") || strings.Contains(errMsg, "ACL not found") {
				code = 403
			}
			rsp.WriteHeader(code)
			rsp.Write([]byte(errMsg))
			return
		}

		// return json object
		if obj != nil {
			var buf []byte
			buf, err = json.Marshal(obj)
			if err != nil {
				goto HAS_ERR
			}
			rsp.Header().Set("Content-Type", "application/json")
			rsp.Write(buf)
		}
	}
	return f
}

func (self *OrderHttpSrv) Index(rsp http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		rsp.WriteHeader(404)
		return
	}
	rsp.Write([]byte("wx web service."))
}

func (self *OrderHttpSrv) decodeBody(req *http.Request, out interface{}, cb func(interface{}) error) error {
	var raw interface{}
	dec := json.NewDecoder(req.Body)
	if err := dec.Decode(&raw); err != nil {
		return err
	}

	if cb != nil {
		if err := cb(raw); err != nil {
			return err
		}
	}

	return mapstructure.Decode(raw, out)
}

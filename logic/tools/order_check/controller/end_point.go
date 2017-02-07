package controller

import (
	"encoding/json"
	"net/http"
	
	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/tools/order_check/fx_models"
)

func (self *OrderHttpSrv) TaobaoOrder(rsp http.ResponseWriter, req *http.Request) (interface{}, error) {
	var request []fx_models.TaobaoOrder
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		logrus.Errorf("TaobaoOrder json decode error: %v", err)
		return nil, err
	}
	
	logrus.Debugf("taobao order request: %v", request)
	
	response := OrderResponse{Code: ORDER_RESPONSE_OK}
	for _, v := range request {
		err := self.orderManager.TaobaoOrder(&v)
		if err != nil {
			logrus.Errorf("taobao order check error: %v", err)
		}
	}
	
	return response, nil
}

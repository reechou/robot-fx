package ext

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/tools/order_check/config"
)

const (
	URI_WX_ROBOT_SEND_MSG = "/sendmsgs"
)

type WxRobotExt struct {
	cfg    *config.Config
	client *http.Client
}

func NewWxRobotExt(cfg *config.Config) *WxRobotExt {
	wre := &WxRobotExt{
		cfg:    cfg,
		client: &http.Client{},
	}

	return wre
}

func (we *WxRobotExt) SendMsg(info *SendMsgInfo) error {
	u := "http://" + we.cfg.WxRobotSrv.Host + URI_WX_ROBOT_SEND_MSG
	body, err := json.Marshal(info)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", u, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	rsp, err := we.client.Do(httpReq)
	defer func() {
		if rsp != nil {
			rsp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var response SendMsgResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		logrus.Errorf("goods search json decode error: %s", string(rspBody))
		return err
	}
	if response.Code != 0 {

	}

	return nil
}

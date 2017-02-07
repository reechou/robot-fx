package ext
//
//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"strings"
//
//	"github.com/Sirupsen/logrus"
//	"github.com/reechou/robot-fx/logic/tools/order_check/config"
//)
//
//const (
//	URI_WX_ROBOT_SEND_MSG = "/index.php?r=sreach/sreach"
//)
//
//type WxRobotExt struct {
//	cfg    *config.Config
//	client *http.Client
//}
//
//func NewWxRobotExt(cfg *config.Config) *WxRobotExt {
//	wre := &WxRobotExt{
//		cfg:    cfg,
//		client: &http.Client{},
//	}
//
//	return wre
//}
//
//func (we *WxRobotExt) SendMsg(info *GoodsSearchReq) (*GoodsSearchData, error) {
//	u := "http://" + we.cfg.DuobbManagerSrv.Host + URI_WX_ROBOT_SEND_MSG
//	body, err := json.Marshal(info)
//	if err != nil {
//		return nil, err
//	}
//	httpReq, err := http.NewRequest("POST", u, strings.NewReader(string(body)))
//	if err != nil {
//		return nil, err
//	}
//	httpReq.Header.Set("Content-Type", "application/json")
//
//	rsp, err := we.client.Do(httpReq)
//	defer func() {
//		if rsp != nil {
//			rsp.Body.Close()
//		}
//	}()
//	if err != nil {
//		return nil, err
//	}
//	rspBody, err := ioutil.ReadAll(rsp.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	var response GoodsSearchRsp
//	err = json.Unmarshal(rspBody, &response)
//	if err != nil {
//		logrus.Errorf("goods search json decode error: %s", string(rspBody))
//		return nil, err
//	}
//	if response.Code != DUOBB_MANAGER_RESPONSE_OK {
//		if response.Code == DUOBB_MANAGER_GOODS_SEARCH_NO_DISCOUNT {
//			return nil, ERR_DUOBB_GOODS_SEARCH_NO_DISCOUNT
//		}
//		logrus.Errorf("goods search error: %v", response)
//		return nil, fmt.Errorf("goods search error: %v", response)
//	}
//
//	return &response.Data, nil
//}

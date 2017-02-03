package ext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/config"
)

const (
	URI_DUOBB_MANAGER_SEARCH = "/index.php?r=sreach/sreach"
)

type DuobbManagerExt struct {
	cfg    *config.Config
	client *http.Client
}

func NewDuobbManagerExt(cfg *config.Config) *DuobbManagerExt {
	dme := &DuobbManagerExt{
		cfg:    cfg,
		client: &http.Client{},
	}

	return dme
}

func (we *DuobbManagerExt) GoodsSearch(info *GoodsSearchReq) (*GoodsSearchData, error) {
	u := "http://" + we.cfg.DuobbManagerSrv.Host + URI_DUOBB_MANAGER_SEARCH
	body, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", u, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	rsp, err := we.client.Do(httpReq)
	defer func() {
		if rsp != nil {
			rsp.Body.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	var response GoodsSearchRsp
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		logrus.Errorf("goods search json decode error: %s", string(rspBody))
		return nil, err
	}
	if response.Code != WECHAT_RESPONSE_OK {
		logrus.Errorf("goods search error: %v", response)
		return nil, fmt.Errorf("goods search error: %v", response)
	}

	return &response.Data, nil
}

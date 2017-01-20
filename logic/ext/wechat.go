package ext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/config"
)

const (
	WECHAT_WITHDRAWAL = "/index.php?r=weixinpay/pay"
	WECHAT_SENDMSG    = "/index.php?r=weixintem/send"
)

type WechatExt struct {
	cfg    *config.Config
	client *http.Client

	stop      chan struct{}
	done      chan struct{}
	wxMsgChan chan *WeixinMsgSendReq
}

func NewWechatExt(cfg *config.Config) *WechatExt {
	we := &WechatExt{
		cfg:       cfg,
		client:    &http.Client{},
		stop:      make(chan struct{}),
		done:      make(chan struct{}),
		wxMsgChan: make(chan *WeixinMsgSendReq, WX_SEND_MSG_CHAN_LEN),
	}
	go we.run()

	return we
}

func (we *WechatExt) Stop() {
	close(we.stop)
	<-we.done
}

func (we *WechatExt) run() {
	logrus.Debugf("wechat ext start run.")
	for {
		select {
		case msg := <-we.wxMsgChan:
			we.WxSendMsg(msg)
		case <-we.stop:
			close(we.done)
			return
		}
	}
}

func (we *WechatExt) AsyncWxSendMsg(msg *WeixinMsgSendReq) {
	logrus.Debugf("async wx send msg[%v]", msg)
	select {
	case we.wxMsgChan <- msg:
	case <-time.After(1 * time.Second):
		logrus.Errorf("async wx send msg timeout.")
		return
	}
}

func (we *WechatExt) Withdrawal(info *WithdrawalReq) error {
	u := we.cfg.WechatExtInfo.HostURL + WECHAT_WITHDRAWAL
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

	var response WechatResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		return err
	}
	if response.Code != WECHAT_RESPONSE_OK {
		logrus.Errorf("wechat withdrawal error: %v", response)
		return fmt.Errorf("wechat withdrawal error: %v", response)
	}

	return nil
}

func (we *WechatExt) WxSendMsg(info *WeixinMsgSendReq) error {
	u := we.cfg.WechatExtInfo.HostURL + WECHAT_SENDMSG
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

	var response WechatResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		return err
	}
	if response.Code != WECHAT_RESPONSE_OK {
		logrus.Errorf("wechat send msg error: %v", response)
		return fmt.Errorf("wechat send msg error: %s", response.Msg)
	}

	return nil
}

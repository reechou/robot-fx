package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	rpcjson "github.com/gorilla/rpc/json"
	"github.com/reechou/duobb_proto"
)

var (
	HostName string
	IP       string
)

func (self *Daemon) CreateAlimamaAdzone(robot, adzone, ali string) (string, string, string, error) {
	// get cookie
	cookie := self.getAlimamaCookie(ali)
	if cookie == "" {
		return "", "", "", fmt.Errorf("get ali cookie == nil")
	}
	// parse tb_token
	var tbToken string
	cs := strings.Split(cookie, ";")
	for _, v := range cs {
		if strings.Contains(v, "_tb_token_") {
			tbToken = strings.Replace(v, " ", "", -1)
			break
		}
	}
	// get guide list
	memberId, guideId, err := self.getAlimamaGuideList(robot, tbToken, cookie)
	if err != nil {
		return "", "", "", err
	}
	if guideId == "" {
		err = self.addAlimamaGuide(robot, tbToken, cookie)
		if err != nil {
			return "", "", "", err
		}
		memberId, guideId, err = self.getAlimamaGuideList(robot, tbToken, cookie)
		if err != nil {
			return "", "", "", err
		}
	}
	adzoneId, err := self.createAlimamaAdzone(adzone, guideId, tbToken, cookie)
	if err != nil {
		return "", "", "", err
	}

	return memberId, guideId, adzoneId, nil
}

func (self *Daemon) createAlimamaAdzone(adzone, siteId, tbToken, cookie string) (string, error) {
	values := url.Values{}
	values.Add("tag", "29")
	values.Add("gcid", "8")
	values.Add("siteid", siteId)
	values.Add("selectact", "add")
	values.Add("newadzonename", adzone)
	values.Add("t", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	values.Add("_tb_token_", strings.Replace(tbToken, "_tb_token_=", "", -1))
	values.Add("pvid", fmt.Sprintf("10_%s_557_%s", IP, strconv.Itoa(int(time.Now().UnixNano()/1000000))))

	req, err := http.NewRequest("POST", ALIMAMA_ADZONE_CREATE, bytes.NewBufferString(values.Encode()))
	if err != nil {
		logrus.Errorf("createAlimamaAdzone http new request error: %v", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "pub.alimama.com")
	req.Header.Set("Referer", "http://pub.alimama.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36")
	req.Header.Set("Cookie", cookie)
	resp, err := self.client.Do(req)
	if err != nil {
		logrus.Errorf("createAlimamaAdzone http do request error: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("createAlimamaAdzone ioutil ReadAll error: %v", err)
		return "", err
	}
	var alimamaRsp AlimamaAdzoneCreateRsp
	err = json.Unmarshal(rspBody, &alimamaRsp)
	if err != nil {
		logrus.Errorf("createAlimamaAdzone json decode error: %v [%s]", err, string(rspBody))
		return "", err
	}
	if alimamaRsp.OK && alimamaRsp.Info.OK {
		return strconv.Itoa(alimamaRsp.Data.AdzoneId), nil
	}
	return "", fmt.Errorf("alimama add guide rsp not ok[%v].", alimamaRsp)
}

func (self *Daemon) getAlimamaGuideList(robot, tbToken, cookie string) (string, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(ALIMAMA_GET_GUIDE_LIST, strconv.Itoa(int(time.Now().UnixNano()/1000000)), tbToken), nil)
	if err != nil {
		logrus.Errorf("getAlimamaGuideList http new request error: %v", err)
		return "", "", err
	}
	req.Header.Set("Host", "pub.alimama.com")
	req.Header.Set("Referer", "http://pub.alimama.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36")
	req.Header.Set("Cookie", cookie)
	resp, err := self.client.Do(req)
	if err != nil {
		logrus.Errorf("getAlimamaGuideList http do request error: %v", err)
		return "", "", err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("getAlimamaGuideList ioutil ReadAll error: %v", err)
		return "", "", err
	}
	var alimamaRsp AlimamaGetGuideListRsp
	err = json.Unmarshal(rspBody, &alimamaRsp)
	if err != nil {
		logrus.Errorf("getAlimamaGuideList json decode error: %v [%s]", err, string(rspBody))
		return "", "", err
	}
	logrus.Debugf("getAlimamaGuideList: %v", alimamaRsp)
	var memberId string
	robotGuide := fmt.Sprintf(ROBOT_DEFAULT_GUIDE, robot)
	if alimamaRsp.OK && alimamaRsp.Info.OK {
		for _, v := range alimamaRsp.Data.GuideList {
			if v.Name == robotGuide {
				return strconv.Itoa(v.MemberId), strconv.Itoa(v.GuideId), nil
			}
			if memberId == "" {
				memberId = strconv.Itoa(v.MemberId)
			}
		}
		return memberId, "", nil
	}
	return "", "", fmt.Errorf("alimama get guide list rsp not ok.")
}

func (self *Daemon) addAlimamaGuide(robot, tbToken, cookie string) error {
	values := url.Values{}
	values.Add("name", fmt.Sprintf(ROBOT_DEFAULT_GUIDE, robot))
	values.Add("categoryId", "14")
	values.Add("account1", "ReeZhou")
	values.Add("pvid", "")
	values.Add("t", strconv.Itoa(int(time.Now().UnixNano()/1000000)))
	values.Add("_tb_token_", strings.Replace(tbToken, "_tb_token_=", "", -1))

	req, err := http.NewRequest("POST", ALIMAMA_GUIDE_ADD, bytes.NewBufferString(values.Encode()))
	if err != nil {
		logrus.Errorf("addAlimamaGuide http new request error: %v", err)
		return err
	}
	logrus.Debug(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Host", "pub.alimama.com")
	req.Header.Set("Referer", "http://pub.alimama.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36")
	req.Header.Set("Cookie", cookie)
	resp, err := self.client.Do(req)
	if err != nil {
		logrus.Errorf("addAlimamaGuide http do request error: %v", err)
		return err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("addAlimamaGuide ioutil ReadAll error: %v", err)
		return err
	}
	var alimamaRsp AlimamaAddGuideRsp
	err = json.Unmarshal(rspBody, &alimamaRsp)
	if err != nil {
		logrus.Errorf("addAlimamaGuide json decode error: %v", err)
		return err
	}
	if alimamaRsp.OK && alimamaRsp.Info.OK {
		return nil
	}
	return fmt.Errorf("alimama add guide rsp not ok[%v].", alimamaRsp)
}

func (self *Daemon) getAlimamaSelfAdzoneList(tbToken, cookie string) (string, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(ALIMAMA_GET_SELF_ADZONE_LIST, strconv.Itoa(int(time.Now().UnixNano()/1000000)), tbToken, IP, strconv.Itoa(int(time.Now().UnixNano()/1000000))), nil)
	if err != nil {
		logrus.Errorf("getAlimamaSelfAdzoneList http new request error: %v", err)
		return "", "", err
	}
	req.Header.Set("Host", "pub.alimama.com")
	req.Header.Set("Referer", "http://pub.alimama.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.95 Safari/537.36")
	req.Header.Set("Cookie", cookie)
	resp, err := self.client.Do(req)
	if err != nil {
		logrus.Errorf("getAlimamaSelfAdzoneList http do request error: %v", err)
		return "", "", err
	}
	defer resp.Body.Close()
	rspBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("getAlimamaSelfAdzoneList ioutil ReadAll error: %v", err)
		return "", "", err
	}
	var alimamaRsp AlimamaGetSelfAdzoneRsp
	err = json.Unmarshal(rspBody, &alimamaRsp)
	if err != nil {
		logrus.Errorf("getAlimamaSelfAdzoneList json decode error: %v", err)
		return "", "", err
	}
	var memberId string
	if alimamaRsp.OK == true {
		for _, v := range alimamaRsp.Data.OtherList {
			if v.Name == ROBOT_DEFAULT_GUIDE {
				return strconv.Itoa(v.MemberId), strconv.Itoa(v.SiteId), nil
			}
			if memberId == "" {
				memberId = strconv.Itoa(v.MemberId)
			}
		}
		return memberId, "", nil
	}
	return "", "", fmt.Errorf("alimama rsp not ok.")
}

func (self *Daemon) getAlimamaCookie(ali string) string {
	request := map[string]string{"ali": ali}
	url := "http://" + self.cfg.DuobbSrv.Host + "/rpc"
	message, err := rpcjson.EncodeClientRequest(DUOBB_GET_ALIMAMA_COOKIE_METHOD, request)
	if err != nil {
		logrus.Errorf("json encode client request error: %v", err)
		return ""
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		logrus.Errorf("http new request error: %v", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("http do request error: %v", err)
		return ""
	}
	defer resp.Body.Close()
	var result duobb_proto.Response
	err = rpcjson.DecodeClientResponse(resp.Body, &result)
	if err != nil {
		logrus.Errorf("json decode client response error: %v", err)
		return ""
	}
	if result.Code == duobb_proto.DUOBB_RSP_SUCCESS {
		dataMap := result.Data.(map[string]interface{})
		if dataMap != nil {
			cookie, ok := dataMap["cookie"]
			if ok {
				return cookie.(string)
			}
		}
	}

	return ""
}

func init() {
	HostName = GetHostName()
	IP = GetLocalIP()
}

func GetHostName() string {
	hostName, err := os.Hostname()
	if err != nil {
		logrus.Errorf("GetHostName error:", err.Error())
		return ""
	}
	return hostName
}

func GetLocalIP() string {
	ipAddress, err := net.ResolveIPAddr("ip", HostName)
	if err != nil {
		logrus.Errorf("GetLocalIP error:", err.Error())
		return ""
	}
	return ipAddress.String()
}

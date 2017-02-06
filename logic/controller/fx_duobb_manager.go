package controller

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/ext"
	"github.com/reechou/robot-fx/logic/models"
)

func (daemon *Daemon) TaobaoGoodsSearch(robot, info string, account *models.FxAccount) (*ext.GoodsSearchData, error) {
	robotAlimama := &models.FxRobotAlimama{
		RobotWX: robot,
	}
	has, err := models.GetRobotAlimama(robotAlimama)
	if err != nil {
		logrus.Errorf("get robot alimama error: %v", err)
		return nil, err
	}
	if !has {
		logrus.Errorf("get robot alimama cannot found alimama from robot[%s]", robot)
		return nil, fmt.Errorf("get robot alimama cannot found alimama from robot[%s]", robot)
	}
	pid := "mm_" + account.MemberId + "_" + account.GuideId + "_" + account.AdzoneId
	if strings.Contains(info, "http") {
		reg := regexp.MustCompile(`(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
		url := reg.FindString(info)
		logrus.Debugf("goods search url: %s", url)
		searchReq := &ext.GoodsSearchReq{
			Alimama:  robotAlimama.Alimama,
			Pid:      pid,
			AdzoneId: account.AdzoneId,
			SiteId:   account.GuideId,
			Url:      url,
		}
		data, err := daemon.dme.GoodsSearch(searchReq)
		if err != nil {
			if err == ext.ERR_DUOBB_GOODS_SEARCH_NO_DISCOUNT {
				return nil, err
			}
			logrus.Errorf("goods search error: %v", err)
			return nil, err
		}
		logrus.Debugf("data: %v", data)
		return data, nil
	}

	return nil, nil
}

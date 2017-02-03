package fx

import (
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/models"
	"strconv"
)

func (fxr *FXRouter) getReqAccount(req *ReceiveMsgInfo) (*models.FxAccount, error) {
	account := &models.FxAccount{RobotWx: req.BaseInfo.WechatNick, UserName: req.BaseInfo.FromUserName}
	has, err := models.GetFxAccountFromUserName(account)
	if err != nil {
		logrus.Errorf("get fx account from username[%v] erorr: %v", account, err)
		return nil, err
	}
	if has {
		if req.BaseInfo.FromNickName != account.Name && req.BaseInfo.FromNickName != "" {
			account.UnionId = fxr.createRobotUnionId(req)
			account.Name = req.BaseInfo.FromNickName
			err = models.UpdateFxAccountName(account)
			if err != nil {
				logrus.Errorf("update fx account name[%v] error: %v", account, err)
				return nil, err
			}
		}
		return account, nil
	}
	logrus.Debugf("cannot found username[%v]", account)
	unionId := fxr.createRobotUnionId(req)
	account = &models.FxAccount{
		UnionId: unionId,
	}
	has, err = models.GetFxAccount(account)
	if err != nil {
		logrus.Errorf("get fx account from unionid[%v] erorr: %v", account, err)
		return nil, err
	}
	if !has {
		logrus.Errorf("cannot found this req[%v] account.", req)
		return nil, REQ_ACCOUNT_GET_NONE
	}
	// update username
	account.UserName = req.BaseInfo.FromUserName
	err = models.UpdateFxAccountUserName(account)
	if err != nil {
		logrus.Errorf("update account[%v] username erorr: %v", account, err)
	}
	return account, nil
}

func (fxr *FXRouter) robotSign(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	a, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	fxAccount := &models.FxAccount{
		UnionId: a.UnionId,
	}
	affected, signScore, err := fxr.backend.UpdateFxAccountSignTime(fxAccount)
	if err != nil {
		logrus.Errorf("Error update fx sign time: %v", err)
		return err
	} else {
		sendMsg := SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
		}
		if affected == 0 {
			sendMsg.Msg = CALLBACK_SIGN_FAILED
		} else {
			sendMsg.Msg = fmt.Sprintf(CALLBACK_SIGN_SUCCESS, signScore)
		}
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, sendMsg)
		return nil
	}

	return nil
}

func (fxr *FXRouter) robotUserInfo(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	a, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	withdrawal, err := fxr.backend.GetWithdrawalRecordSum(a.UnionId)
	if err != nil {
		logrus.Errorf("get fx withdrawal record sum error: %v", err)
		return err
	}
	orderCount, err := fxr.backend.GetFxOrderListCount(a.UnionId)
	if err != nil {
		logrus.Errorf("get fx order list count error: %v", err)
		return err
	}
	waitOrderCount, err := fxr.backend.GetFxOrderWaitSettlementRecordListCount(a.UnionId)
	if err != nil {
		logrus.Errorf("get fx wait order list count error: %v", err)
		return err
	}
	waitSum, err := fxr.backend.GetFxOrderWaitSettlementRecordSum(a.ID)
	if err != nil {
		logrus.Errorf("get fx wait order sum error: %v", err)
		return err
	}
	rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		WechatNick: req.BaseInfo.WechatNick,
		ChatType:   CHAT_TYPE_PEOPLE,
		NickName:   req.BaseInfo.FromNickName,
		UserName:   req.BaseInfo.FromUserName,
		MsgType:    MSG_TYPE_TEXT,
		Msg: fmt.Sprintf(CALLBACK_USER_INFO_SUCCESS, req.BaseInfo.FromNickName, int(a.CanWithdrawals), int(a.AllScore),
			int(withdrawal), int(orderCount), int(waitOrderCount), int(waitSum)),
	})

	return nil
}

func (fxr *FXRouter) robotGetLowerPeople(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	a, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	count, err := fxr.backend.GetLowerPeopleCount(a.WechatUnionId)
	if err != nil {
		logrus.Errorf("get lower people count error: %v", err)
		return err
	}
	list, err := fxr.backend.GetLowerPeopleList(a.WechatUnionId, 0, 20)
	if err != nil {
		logrus.Errorf("get lower people list error: %v", err)
		return err
	}
	msg := fmt.Sprintf(CALLBACK_LOWER_PEOPLE_SUCCESS, count)
	for i, v := range list {
		msg += "\n" + strconv.Itoa(i+1) + ". " + v.Name
	}
	rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		WechatNick: req.BaseInfo.WechatNick,
		ChatType:   CHAT_TYPE_PEOPLE,
		NickName:   req.BaseInfo.FromNickName,
		UserName:   req.BaseInfo.FromUserName,
		MsgType:    MSG_TYPE_TEXT,
		Msg:        msg,
	})

	return nil
}

func (fxr *FXRouter) robotBindWechat(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	wechat := strings.Replace(req.Msg, KEYWORD_BIND_WECHAT, "", -1)

	a, err := fxr.getReqAccount(req)
	if err != nil {
		if err != REQ_ACCOUNT_GET_NONE {
			return err
		}
		// create fx account
		unionId := fxr.createRobotUnionId(req)
		a = &models.FxAccount{
			UnionId:  unionId,
			RobotWx:  req.BaseInfo.WechatNick,
			UserName: req.BaseInfo.FromUserName,
			Name:     req.BaseInfo.FromNickName,
		}
		err = models.CreateFxAccount(a)
		if err != nil {
			logrus.Errorf("create fx account error: %v", err)
			return err
		}
	}
	if a.WechatUnionId != "" {
		logrus.Debugf("fx account[%s] has exist wechat.", a.UnionId)
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
			Msg:        fmt.Sprintf(CALLBACK_BIND_WECHAT_FAILED, req.BaseInfo.FromNickName),
		})
		return nil
	}

	wechatUnionId := fxr.createWechatUnionId(req.BaseInfo.WechatNick, wechat)
	// create alimama guide
	robotAlimama := &models.FxRobotAlimama{
		RobotWX: req.BaseInfo.WechatNick,
	}
	has, err := models.GetRobotAlimama(robotAlimama)
	if err != nil {
		logrus.Errorf("get fx robot[%s] alimama error: %v", req.BaseInfo.WechatNick, err)
		return err
	}
	if !has {
		logrus.Errorf("cannot found robot[%s] alimama", req.BaseInfo.WechatNick)
		return fmt.Errorf("cannot found robot[%s] alimama", req.BaseInfo.WechatNick)
	}
	memberId, guideId, adzoneId, err := fxr.backend.CreateAlimamaAdzone(req.BaseInfo.WechatNick, wechatUnionId, robotAlimama.Alimama)
	if err != nil {
		logrus.Errorf("create alimama adzone error: %v", err)
		return err
	}
	fxAccount := &models.FxAccount{
		UnionId:       a.UnionId,
		WechatUnionId: wechatUnionId,
		MemberId:      memberId,
		GuideId:       guideId,
		AdzoneId:      adzoneId,
	}
	err = models.UpdateFxAccountWechat(fxAccount)
	if err != nil {
		logrus.Errorf("update fx account wechat error: %v", err)
		return err
	}
	rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		WechatNick: req.BaseInfo.WechatNick,
		ChatType:   CHAT_TYPE_PEOPLE,
		NickName:   req.BaseInfo.FromNickName,
		UserName:   req.BaseInfo.FromUserName,
		MsgType:    MSG_TYPE_TEXT,
		Msg:        fmt.Sprintf(CALLBACK_BIND_WECHAT_SUCCESS, req.BaseInfo.FromNickName, wechat),
	})

	return nil
}

func (fxr *FXRouter) robotOrderList(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	a, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	list, err := fxr.backend.GetFxAllOrderList(a.UnionId, 0, 10)
	if err != nil {
		logrus.Errorf("get fx all order list error: %v", err)
		return err
	}
	msg := CALLBACK_ORDER_LIST_SUCCESS
	for _, v := range list {
		msg += "\n" + fmt.Sprintf("%s**** %s", v.OrderId[:4], time.Unix(v.UpdatedAt, 0).Format("2006-01-02"))
		if v.Status == FX_ORDER_SETTLEMENT {
			msg += " 已结算"
		}
		msg += fmt.Sprintf("返 %d", int(v.ReturnMoney))
	}
	rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		WechatNick: req.BaseInfo.WechatNick,
		ChatType:   CHAT_TYPE_PEOPLE,
		NickName:   req.BaseInfo.FromNickName,
		UserName:   req.BaseInfo.FromUserName,
		MsgType:    MSG_TYPE_TEXT,
		Msg:        msg,
	})

	return nil
}

func (fxr *FXRouter) robotGoodsSearch(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	a, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	data, err := fxr.backend.TaobaoGoodsSearch(req.BaseInfo.WechatNick, req.Msg, a)
	if err != nil {
		return err
	}
	var rate float32
	if len(fxr.cfg.LevelPer) != 0 {
		rate = float32(fxr.cfg.LevelPer[0])
	} else {
		rate = DEFAULT_RETURN_RATE
	}
	returnMoney := data.EndPrice * data.RlRate * rate / 10000.0

	if data.Amount != 0.0 {
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
			Msg: fmt.Sprintf(CALLBACK_GOODS_SEARCH_SUCCESS, a.Name, data.Title, data.ZkPrice, data.EndPrice-returnMoney,
				data.Amount+returnMoney, data.Amount, returnMoney, data.Token),
		})
	} else {
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
			Msg: fmt.Sprintf(CALLBACK_GOODS_SEARCH_NO_QUAN_SUCCESS, a.Name, data.Title, data.ZkPrice, data.EndPrice-returnMoney,
				returnMoney, data.Token),
		})
	}

	return nil
}

func (fxr *FXRouter) robotWithdrawal(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	a, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	wInfo := &models.WithdrawalRecord{
		UnionId:         a.UnionId,
		WithdrawalMoney: a.CanWithdrawals,
	}
	err, ifSystemErr := fxr.backend.CreateWithdrawalRecord(wInfo, a)
	if err != nil {
		if !ifSystemErr {
			rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
				WechatNick: req.BaseInfo.WechatNick,
				ChatType:   CHAT_TYPE_PEOPLE,
				NickName:   req.BaseInfo.FromNickName,
				UserName:   req.BaseInfo.FromUserName,
				MsgType:    MSG_TYPE_TEXT,
				Msg:        err.Error(),
			})
			return nil
		}
		logrus.Errorf("create withdrawal record error: %v", err)
		return err
	}
	withdrawalMoney := wInfo.WithdrawalMoney / float32(fxr.cfg.Score.EnlargeScale)
	rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		WechatNick: req.BaseInfo.WechatNick,
		ChatType:   CHAT_TYPE_PEOPLE,
		NickName:   req.BaseInfo.FromNickName,
		UserName:   req.BaseInfo.FromUserName,
		MsgType:    MSG_TYPE_TEXT,
		Msg: fmt.Sprintf(CALLBACK_WITHDRAWAL_SUCCESS, a.Name, a.CanWithdrawals, withdrawalMoney,
			fxr.cfg.WithdrawalPolicy.MonthWithdrawalTime, fxr.cfg.WithdrawalPolicy.MinimumWithdrawal),
	})
	rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		WechatNick: req.BaseInfo.WechatNick,
		ChatType:   CHAT_TYPE_PEOPLE,
		NickName:   fxr.cfg.WithdrawalPolicy.NotifyPeople,
		MsgType:    MSG_TYPE_TEXT,
		Msg: fmt.Sprintf(CALLBACK_WITHDRAWAL_NOTIFY, req.BaseInfo.WechatNick, a.Name, a.Wechat, withdrawalMoney,
			time.Now().Format("2006-01-02 15:04:05")),
	})

	return nil
}

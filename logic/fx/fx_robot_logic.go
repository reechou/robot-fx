package fx

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/ext"
	"github.com/reechou/robot-fx/logic/models"
)

func (fxr *FXRouter) getReqAccount(req *ReceiveMsgInfo) (*models.FxAccount, *models.FxWxAccount, error) {
	//req.BaseInfo.FromNickName = fxr.filterEmoji(req.BaseInfo.FromNickName)

	account := &models.FxAccount{RobotWx: req.BaseInfo.WechatNick, UserName: req.BaseInfo.FromUserName}
	has, err := models.GetFxAccountFromUserName(account)
	if err != nil {
		logrus.Errorf("get fx account from username[%v] erorr: %v", account, err)
		return nil, nil, err
	}
	if has {
		if req.BaseInfo.FromNickName != account.Name && req.BaseInfo.FromNickName != "" {
			account.UnionId = fxr.createRobotUnionId(req)
			account.Name = req.BaseInfo.FromNickName
			err = models.UpdateFxAccountName(account)
			if err != nil {
				logrus.Errorf("update fx account name[%v] error: %v", account, err)
				return nil, nil, err
			}
			updateFxWxAccount := &models.FxWxAccount{
				WxId: account.WxId,
				Name: req.BaseInfo.FromNickName,
			}
			err = models.UpdateFxWxAccountName(updateFxWxAccount)
			if err != nil {
				logrus.Errorf("update fx wx account name[%v] error: %v", updateFxWxAccount, err)
			}
		}
		fxWxAccount := &models.FxWxAccount{
			WxId: account.WxId,
		}
		has, err = models.GetFxWxAccount(fxWxAccount)
		if err != nil {
			logrus.Errorf("get fx wx account error: %v", err)
			return nil, nil, err
		}
		if !has {
			logrus.Errorf("fx wx account has none this account")
			return nil, nil, REQ_ACCOUNT_GET_NONE
		}
		return account, fxWxAccount, nil
	}
	logrus.Debugf("cannot found username[%v]", account)
	unionId := fxr.createRobotUnionId(req)
	account = &models.FxAccount{
		UnionId: unionId,
	}
	has, err = models.GetFxAccount(account)
	if err != nil {
		logrus.Errorf("get fx account from unionid[%v] erorr: %v", account, err)
		return nil, nil, err
	}
	if !has {
		logrus.Errorf("cannot found this req[%v] account.", req)
		return nil, nil, REQ_ACCOUNT_GET_NONE
	}
	// update username
	account.UserName = req.BaseInfo.FromUserName
	err = models.UpdateFxAccountUserName(account)
	if err != nil {
		logrus.Errorf("update account[%v] username erorr: %v", account, err)
	}
	fxWxAccount := &models.FxWxAccount{
		WxId: account.WxId,
	}
	has, err = models.GetFxWxAccount(fxWxAccount)
	if err != nil {
		logrus.Errorf("get fx wx account error: %v", err)
		return nil, nil, err
	}
	if !has {
		logrus.Errorf("fx wx account has none this account")
		return nil, nil, REQ_ACCOUNT_GET_NONE
	}
	return account, fxWxAccount, nil
}

func (fxr *FXRouter) robotHelp(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		WechatNick: req.BaseInfo.WechatNick,
		ChatType:   CHAT_TYPE_PEOPLE,
		NickName:   req.BaseInfo.FromNickName,
		UserName:   req.BaseInfo.FromUserName,
		MsgType:    MSG_TYPE_TEXT,
		Msg:        CALLBACK_HELP,
	})
	return nil
}

func (fxr *FXRouter) robotSign(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	_, wa, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	affected, err := fxr.backend.UpdateFxWxAccountSignTime(wa)
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
			sendMsg.Msg = fmt.Sprintf(CALLBACK_SIGN_FAILED, req.BaseInfo.FromNickName)
		} else {
			sendMsg.Msg = fmt.Sprintf(CALLBACK_SIGN_SUCCESS, req.BaseInfo.FromNickName, fxr.cfg.Score.SignScore)
		}
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, sendMsg)
		return nil
	}

	return nil
}

func (fxr *FXRouter) robotUserInfo(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	_, wa, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	withdrawal, err := fxr.backend.GetWithdrawalRecordSum(wa.WxId)
	if err != nil {
		logrus.Errorf("get fx withdrawal record sum error: %v", err)
		return err
	}
	orderCount, err := fxr.backend.GetFxOrderListCount(wa.WxId)
	if err != nil {
		logrus.Errorf("get fx order list count error: %v", err)
		return err
	}
	waitOrderCount, err := fxr.backend.GetFxOrderWaitSettlementRecordListCount(wa.WxId)
	if err != nil {
		logrus.Errorf("get fx wait order list count error: %v", err)
		return err
	}
	waitSum, err := fxr.backend.GetFxOrderWaitSettlementRecordSum(wa.ID)
	if err != nil {
		logrus.Errorf("get fx wait order sum from account[%v] error: %v", wa, err)
		return err
	}
	rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		WechatNick: req.BaseInfo.WechatNick,
		ChatType:   CHAT_TYPE_PEOPLE,
		NickName:   req.BaseInfo.FromNickName,
		UserName:   req.BaseInfo.FromUserName,
		MsgType:    MSG_TYPE_TEXT,
		Msg: fmt.Sprintf(CALLBACK_USER_INFO_SUCCESS, req.BaseInfo.FromNickName, int(wa.CanWithdrawals), int(wa.AllScore),
			int(withdrawal), int(orderCount), int(waitOrderCount), int(waitSum)),
	})

	return nil
}

func (fxr *FXRouter) robotGetLowerPeople(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	_, wa, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	count, err := models.GetFxWxLowerPeopleCount(wa.WxId)
	if err != nil {
		logrus.Errorf("get lower people count error: %v", err)
		return err
	}
	list, err := models.GetFxWxLowerPeople(wa.WxId, 0, 20)
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

	a, _, err := fxr.getReqAccount(req)
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
	_, wa, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	list, err := fxr.backend.GetFxAllOrderList(wa.WxId, 0, 10)
	if err != nil {
		logrus.Errorf("get fx all order list error: %v", err)
		return err
	}
	msg := CALLBACK_ORDER_LIST_SUCCESS
	for i, v := range list {
		nameRune := []rune(v.OrderName)
		msg += "\n" + fmt.Sprintf("(%d) %s**** %s**** %s", i+1, v.OrderId[:4], string(nameRune[:6]), time.Unix(v.UpdatedAt, 0).Format("2006-01-02"))
		if v.Status == FX_ORDER_SETTLEMENT {
			msg += " 已结算"
		}
		msg += fmt.Sprintf(" 约返%d积分", int(v.ReturnMoney*float32(fxr.cfg.Score.EnlargeScale*fxr.cfg.SettlementCommission.LevelPer[0]/100)))
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
	a, _, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	data, err := fxr.backend.TaobaoGoodsSearch(req.BaseInfo.WechatNick, req.Msg, a)
	if err != nil {
		if err == ext.ERR_DUOBB_GOODS_SEARCH_NO_DISCOUNT {
			rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
				WechatNick: req.BaseInfo.WechatNick,
				ChatType:   CHAT_TYPE_PEOPLE,
				NickName:   req.BaseInfo.FromNickName,
				UserName:   req.BaseInfo.FromUserName,
				MsgType:    MSG_TYPE_TEXT,
				Msg:        fmt.Sprintf(CALLBACK_GOODS_NO_DISCOUNT, a.Name),
			})
			return nil
		}
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
				data.Amount+returnMoney, data.Amount, returnMoney),
		})
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
			Msg:        fmt.Sprintf(CALLBACK_PLACE_ORDER, data.Token),
		})
	} else {
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
			Msg:        fmt.Sprintf(CALLBACK_GOODS_SEARCH_NO_QUAN_SUCCESS, a.Name, data.Title, data.ZkPrice, data.EndPrice-returnMoney, returnMoney),
		})
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
			Msg:        fmt.Sprintf(CALLBACK_PLACE_ORDER, data.Token),
		})
	}

	return nil
}

func (fxr *FXRouter) robotWithdrawal(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	a, wa, err := fxr.getReqAccount(req)
	if err != nil {
		logrus.Errorf("get req account error: %v", err)
		return err
	}
	wInfo := &models.WithdrawalRecord{
		UnionId:         wa.WxId,
		WithdrawalMoney: wa.CanWithdrawals,
	}
	err, ifSystemErr := fxr.backend.CreateWithdrawalRecord(wInfo, a, wa)
	if err != nil {
		if !ifSystemErr {
			rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
				WechatNick: req.BaseInfo.WechatNick,
				ChatType:   CHAT_TYPE_PEOPLE,
				NickName:   req.BaseInfo.FromNickName,
				UserName:   req.BaseInfo.FromUserName,
				MsgType:    MSG_TYPE_TEXT,
				Msg:        fmt.Sprintf("%s\n"+CALLBACK_WITHDRAWAL_POLICY, err.Error(), fxr.cfg.WithdrawalPolicy.MonthWithdrawalTime, fxr.cfg.WithdrawalPolicy.MinimumWithdrawal),
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
		Msg: fmt.Sprintf(CALLBACK_WITHDRAWAL_SUCCESS, wa.Name, int64(wa.CanWithdrawals), withdrawalMoney,
			fxr.cfg.WithdrawalPolicy.MonthWithdrawalTime, fxr.cfg.WithdrawalPolicy.MinimumWithdrawal),
	})
	for _, v := range fxr.cfg.WithdrawalPolicy.NotifyPeople {
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   v,
			MsgType:    MSG_TYPE_TEXT,
			Msg: fmt.Sprintf(CALLBACK_WITHDRAWAL_NOTIFY, req.BaseInfo.WechatNick, wa.Name, wa.Wechat, withdrawalMoney,
				time.Now().Format("2006-01-02 15:04:05")),
		})
	}

	return nil
}

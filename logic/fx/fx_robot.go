package fx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/models"
	"github.com/reechou/robot-fx/utils"
	"golang.org/x/net/context"
)

func (fxr *FXRouter) createRobotAlimama(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &CreateRobotAlimamaReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	robotAlimama := &models.FxRobotAlimama{
		RobotWX: req.RobotWX,
		Alimama: req.Alimama,
	}
	err := models.CreateFxRobotAlimama(robotAlimama)
	if err != nil {
		logrus.Errorf("create robot alimama error: %v", err)
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) robotCall(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &ReceiveMsgInfo{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &CallbackMsgInfo{}
	rsp.RetResponse.Code = RspCodeOK
	rsp.BaseInfo = req.BaseInfo

	switch req.BaseInfo.ReceiveEvent {
	case RECEIVE_EVENT_ADD_FRIEND:
		err := fxr.robotAddFriend(req, rsp)
		if err != nil {
			rsp.RetResponse.Code = RspCodeErr
		}
	case RECEIVE_EVENT_MSG:
		err := fxr.robotHandleMsg(req, rsp)
		if err != nil {
			rsp.RetResponse.Code = RspCodeErr
		}
	default:
		logrus.Errorf("robot unknown receive event: %s", req.BaseInfo.ReceiveEvent)
		return nil
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) robotHandleMsg(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	if strings.Contains(req.Msg, KEYWORD_HELP) || req.Msg == KEYWORD_HELP_ID {
		return fxr.robotHelp(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_USER_INFO) || req.Msg == KEYWORD_USER_INFO_ID {
		return fxr.robotUserInfo(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_ORDER_INFO) || req.Msg == KEYWORD_ORDER_INFO_ID {
		return fxr.robotOrderList(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_BIND_WECHAT) {
		return fxr.robotBindWechat(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_SIGN) || req.Msg == KEYWORD_SIGN_ID {
		return fxr.robotSign(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_LOWER_PEOPLE) || req.Msg == KEYWORD_LOWER_PEOPLE_ID {
		return fxr.robotGetLowerPeople(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_WITHDRAWAL) || req.Msg == KEYWORD_WITHDRAWAL_ID {
		return fxr.robotWithdrawal(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_GOODS_SEARCH_URL) {
		return fxr.robotGoodsSearch(req, rsp)
	}

	return nil
}

func (fxr *FXRouter) robotAddFriend(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	fxWxAccount := &models.FxWxAccount{
		WxId: req.AddFriend.UserWxid,
	}
	has, err := models.GetFxWxAccount(fxWxAccount)
	if err != nil {
		logrus.Errorf("get fx wx account from wxid[%s] error: %v", req.AddFriend.UserWxid, err)
		return err
	}
	if !has {
		fxWxAccount.Wechat = req.AddFriend.UserWechat
		fxWxAccount.Name = req.BaseInfo.FromNickName
		if req.AddFriend.SourceWechat != "" {
			fxWxAccount.Superior = req.AddFriend.SourceWechat
		}
		err = fxr.backend.CreateFxWxAccount(fxWxAccount)
		if err != nil {
			logrus.Errorf("create fx wx account error: %v", err)
			return err
		}
	}

	account := &models.FxAccount{
		RobotWx:  req.BaseInfo.WechatNick,
		UserName: req.BaseInfo.FromUserName,
	}
	has, err = models.GetFxAccountFromUserName(account)
	if err != nil {
		logrus.Errorf("get fx account from username[%v] error: %v", account, err)
		return err
	}
	if has {
		if req.BaseInfo.FromNickName != account.Name && req.BaseInfo.FromNickName != "" {
			account.UnionId = fxr.createRobotUnionId(req)
			account.Name = req.BaseInfo.FromNickName
			err = models.UpdateFxAccountName(account)
			if err != nil {
				logrus.Errorf("update fx account name[%v] error: %v", account, err)
				return err
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
		return nil
	}

	unionId := fxr.createRobotUnionId(req)
	// check fx account
	has, err = models.GetFxAccount(&models.FxAccount{UnionId: unionId})
	if err != nil {
		logrus.Errorf("get fx account[%s] error: %v", unionId, err)
		return err
	}
	if has {
		logrus.Debugf("fx account[%s] has exist.", unionId)
		// update username
		account.UnionId = unionId
		account.UserName = req.BaseInfo.FromUserName
		err = models.UpdateFxAccountUserName(account)
		if err != nil {
			logrus.Errorf("update account[%v] username erorr: %v", account, err)
		}
		return nil
	}

	var wechatUnionId, memberId, guideId, adzoneId, superior string
	if req.AddFriend.UserWechat != "" {
		wechatUnionId = fxr.createWechatUnionId(req.BaseInfo.WechatNick, req.AddFriend.UserWxid)
		// create alimama guide
		robotAlimama := &models.FxRobotAlimama{
			RobotWX: req.BaseInfo.WechatNick,
		}
		has, err = models.GetRobotAlimama(robotAlimama)
		if err != nil {
			logrus.Errorf("get fx robot[%s] alimama error: %v", req.BaseInfo.WechatNick, err)
			return err
		}
		if !has {
			logrus.Errorf("cannot found robot[%s] alimama", req.BaseInfo.WechatNick)
			return fmt.Errorf("cannot found robot[%s] alimama", req.BaseInfo.WechatNick)
		}
		memberId, guideId, adzoneId, err = fxr.backend.CreateAlimamaAdzone(req.BaseInfo.WechatNick, wechatUnionId, robotAlimama.Alimama)
		if err != nil {
			logrus.Errorf("create alimama adzone error: %v", err)
			//return err
		}
	}

	if req.AddFriend.SourceWechat != "" && req.AddFriend.SourceNick != req.BaseInfo.WechatNick {
		superior = fxr.createWechatUnionId(req.BaseInfo.WechatNick, req.AddFriend.SourceWechat)
	}
	fxAccount := &models.FxAccount{
		UnionId:       unionId,
		RobotWx:       req.BaseInfo.WechatNick,
		WechatUnionId: wechatUnionId,
		Wechat:        req.AddFriend.UserWechat,
		WxId:          req.AddFriend.UserWxid,
		UserName:      req.BaseInfo.FromUserName,
		Name:          req.BaseInfo.FromNickName,
		Superior:      superior,
		MemberId:      memberId,
		GuideId:       guideId,
		AdzoneId:      adzoneId,
	}
	_, err = fxr.backend.CreateFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("backend create account error: %v", err)
		return err
	}
	if req.AddFriend.UserWechat != "" {
		//rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
		//	WechatNick: req.BaseInfo.WechatNick,
		//	ChatType:   CHAT_TYPE_PEOPLE,
		//	NickName:   req.BaseInfo.FromNickName,
		//	UserName:   req.BaseInfo.FromUserName,
		//	MsgType:    MSG_TYPE_TEXT,
		//	Msg:        CALLBACK_CREATE_ACCOUNT_SUCCESS,
		//})
	} else {
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
			Msg:        CALLBACK_CREATE_ACCOUNT_WITHOUT_WECHAT,
		})
	}
	if fxWxAccount.Superior != "" {
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.AddFriend.SourceNick,
			MsgType:    MSG_TYPE_TEXT,
			Msg:        fmt.Sprintf(CALLBACK_INVITE_SUCCESS, req.AddFriend.UserNick, fxr.cfg.Score.FollowScore),
		})
	}

	return nil
}

func (fxr *FXRouter) createRobotUnionId(req *ReceiveMsgInfo) string {
	return req.BaseInfo.WechatNick + UNION_ID_DELIMITER + req.BaseInfo.FromNickName
}

func (fxr *FXRouter) createWechatUnionId(robotWx, wxid string) string {
	return robotWx + UNION_ID_DELIMITER + wxid
}

func (fxr *FXRouter) filterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

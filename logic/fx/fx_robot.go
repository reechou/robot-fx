package fx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	if strings.Contains(req.Msg, KEYWORD_USER_INFO) {
		return fxr.robotUserInfo(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_ORDER_INFO) {
		return fxr.robotOrderList(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_BIND_WECHAT) {
		return fxr.robotBindWechat(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_SIGN) {
		return fxr.robotSign(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_LOWER_PEOPLE) {
		return fxr.robotGetLowerPeople(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_WITHDRAWAL) {
		return fxr.robotWithdrawal(req, rsp)
	} else if strings.Contains(req.Msg, KEYWORD_GOODS_SEARCH_URL) {
		return fxr.robotGoodsSearch(req, rsp)
	}

	return nil
}

func (fxr *FXRouter) robotAddFriend(req *ReceiveMsgInfo, rsp *CallbackMsgInfo) error {
	account := &models.FxAccount{
		RobotWx:  req.BaseInfo.WechatNick,
		UserName: req.BaseInfo.FromUserName,
	}
	has, err := models.GetFxAccountFromUserName(account)
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
			return err
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
		UserName:      req.BaseInfo.FromUserName,
		Name:          req.BaseInfo.FromNickName,
		Superior:      superior,
		MemberId:      memberId,
		GuideId:       guideId,
		AdzoneId:      adzoneId,
	}
	inviteScore, err := fxr.backend.CreateFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("backend create account error: %v", err)
		return err
	}
	if req.AddFriend.UserWechat != "" {
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.BaseInfo.FromNickName,
			UserName:   req.BaseInfo.FromUserName,
			MsgType:    MSG_TYPE_TEXT,
			Msg:        CALLBACK_CREATE_ACCOUNT_SUCCESS,
		})
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
	if superior != "" {
		rsp.CallbackMsgs = append(rsp.CallbackMsgs, SendBaseInfo{
			WechatNick: req.BaseInfo.WechatNick,
			ChatType:   CHAT_TYPE_PEOPLE,
			NickName:   req.AddFriend.SourceNick,
			MsgType:    MSG_TYPE_TEXT,
			Msg:        fmt.Sprintf(CALLBACK_INVITE_SUCCESS, req.AddFriend.UserNick, inviteScore),
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

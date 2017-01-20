package fx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/logic/models"
	"github.com/reechou/robot-fx/utils"
	"golang.org/x/net/context"
)

func (fxr *FXRouter) getFxAccountHistoryList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxAccountHistoryListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	//fxAccount := &models.FxAccount{
	//	UnionId: req.UnionId,
	//}
	//err := fxr.backend.GetFxAccount(fxAccount)
	//if err != nil {
	//	logrus.Errorf("Error change union_id to account_id error: %v", err)
	//	rsp.Code = RspCodeErr
	//	rsp.Msg = fmt.Sprintf("Error change union_id to account_id error: %v", err)
	//} else {
	type FxAccountHistoryList struct {
		Count int64                     `json:"count"`
		List  []models.FxAccountHistory `json:"list"`
	}
	count, err := fxr.backend.GetFxAccountHistoryListCount(req.UnionId)
	if err != nil {
		logrus.Errorf("Error get fx account history list count: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx account history list count: %v", err)
	} else {
		list, err := fxr.backend.GetFxAccountHistoryList(req.UnionId, req.Offset, req.Num)
		if err != nil {
			logrus.Errorf("Error get fx account history list: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx account history list: %v", err)
		} else {
			var listInfo FxAccountHistoryList
			listInfo.Count = count
			listInfo.List = list
			rsp.Data = listInfo
		}
	}
	//}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxAccountHistoryListByType(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxAccountHistoryListByTypeReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	//fxAccount := &models.FxAccount{
	//	UnionId: req.UnionId,
	//}
	//err := fxr.backend.GetFxAccount(fxAccount)
	//if err != nil {
	//	logrus.Errorf("Error change union_id to account_id error: %v", err)
	//	rsp.Code = RspCodeErr
	//	rsp.Msg = fmt.Sprintf("Error change union_id to account_id error: %v", err)
	//} else {
	type FxAccountHistoryList struct {
		Count int64                     `json:"count"`
		List  []models.FxAccountHistory `json:"list"`
	}
	count, err := fxr.backend.GetFxAccountHistoryListByTypeCount(req.UnionId, req.Type)
	if err != nil {
		logrus.Errorf("Error get fx account history list by type count: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx account history list by type count: %v", err)
	} else {
		list, err := fxr.backend.GetFxAccountHistoryListByType(req.UnionId, req.Type, req.Offset, req.Num)
		if err != nil {
			logrus.Errorf("Error get fx account history list by type: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx account history list by type: %v", err)
		} else {
			var listInfo FxAccountHistoryList
			listInfo.Count = count
			listInfo.List = list
			rsp.Data = listInfo
		}
	}
	//}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

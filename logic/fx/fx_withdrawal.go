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

func (fxr *FXRouter) createFxWithdrawalRecord(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &withdrawalMoneyReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	wInfo := &models.WithdrawalRecord{
		UnionId:         req.UnionId,
		WithdrawalMoney: req.Money,
		OpenId:          req.OpenId,
	}
	err := fxr.backend.CreateWithdrawalRecord(wInfo)
	if err != nil {
		logrus.Errorf("create withdrawal record[%v] error: %v", wInfo, err)
		rsp.Code = RspCodeErr
		rsp.Msg = err.Error()
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxWithdrawalRecordList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getWithdrawalListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	type FxWithdrawalRecordList struct {
		Count int64                     `json:"count"`
		List  []models.WithdrawalRecord `json:"list"`
	}
	count, err := fxr.backend.GetWithdrawalRecordListCount(req.UnionId, req.Status)
	if err != nil {
		logrus.Errorf("Error get fx withdrawal record list count: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx withdrawal record list count: %v", err)
	} else {
		list, err := fxr.backend.GetWithdrawalRecordList(req.UnionId, req.Offset, req.Num, req.Status)
		if err != nil {
			logrus.Errorf("Error get fx withdrawal record list: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx withdrawal record list: %v", err)
		} else {
			var listInfo FxWithdrawalRecordList
			listInfo.Count = count
			listInfo.List = list
			rsp.Data = listInfo
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxWithdrawalRecordSum(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getWithdrawalSumReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	total, err := fxr.backend.GetWithdrawalRecordSum(req.UnionId)
	if err != nil {
		logrus.Errorf("Error get fx withdrawal record sum: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx withdrawal record sum: %v", err)
	} else {
		rsp.Data = total
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxWithdrawalRecordErrorList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getWithdrawalErrorListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	type FxWithdrawalRecordErrorList struct {
		Count int64                          `json:"count"`
		List  []models.WithdrawalRecordError `json:"list"`
	}
	count, err := fxr.backend.GetWithdrawalErrorRecordListCount()
	if err != nil {
		logrus.Errorf("Error get fx withdrawal error msg record list count: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx withdrawal error msg record list count: %v", err)
	} else {
		list, err := fxr.backend.GetWithdrawalErrorRecordList(req.Offset, req.Num)
		if err != nil {
			logrus.Errorf("Error get fx withdrawal error msg record list: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx withdrawal error msg record list: %v", err)
		} else {
			var listInfo FxWithdrawalRecordErrorList
			listInfo.Count = count
			listInfo.List = list
			rsp.Data = listInfo
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxWithdrawalRecordErrorListFromName(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}
	
	req := &getWithdrawalErrorListFromNameReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	
	rsp := &FxResponse{Code: RspCodeOK}
	
	list, err := fxr.backend.GetWithdrawalErrorRecordListFromName(req.Name)
	if err != nil {
		logrus.Errorf("Error get fx withdrawal error msg record list from name: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx withdrawal error msg record list from name: %v", err)
	} else {
		rsp.Data = list
	}
	
	return utils.WriteJSON(w, http.StatusOK, rsp)
}

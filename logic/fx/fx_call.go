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

func (fxr *FXRouter) createFxAccount(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &CreateFxAccountReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId:  req.UnionId,
		Superior: req.Superior,
		Name:     req.Name,
	}
	if _, err := fxr.backend.CreateFxAccount(fxAccount); err != nil {
		logrus.Errorf("Error create fx account: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error create fx account: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) createFxSalesman(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &CreateSalesmanReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	if err := fxr.backend.CreateSalesman(fxAccount); err != nil {
		logrus.Errorf("Error create fx salesman: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error create fx salesman: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) updateFxBaseInfo(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &updateFxBaseInfoReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		Name:  req.Name,
	}
	if err := fxr.backend.UpdateFxAccountBaseInfo(fxAccount); err != nil {
		logrus.Errorf("Error update fx base info: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error update fx base info: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) updateFxStatus(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &updateFxStatusReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
		Status:  req.Status,
	}
	if err := fxr.backend.UpdateFxAccountStatus(fxAccount); err != nil {
		logrus.Errorf("Error update fx status: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error update fx status: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) updateFxSignTime(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &updateFxSignTimeReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	affected, _, err := fxr.backend.UpdateFxAccountSignTime(fxAccount)
	if err != nil {
		logrus.Errorf("Error update fx sign time: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error update fx sign time: %v", err)
	} else {
		if affected == 0 {
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("今天已经签过到了哦!")
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxAccount(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxAccountReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	if err := fxr.backend.GetFxAccount(fxAccount); err != nil {
		logrus.Errorf("Error get fx account: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx account: %v", err)
	} else {
		rsp.Data = fxAccount
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxLowerPeopleList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxLowerPeopleListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	type FxLowerPeopleList struct {
		Count int64              `json:"count"`
		List  []models.FxAccount `json:"list"`
	}
	count, err := fxr.backend.GetLowerPeopleCount(req.UnionId)
	if err != nil {
		logrus.Errorf("Error get fx lower people list count: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx lower people list count: %v", err)
	} else {
		list, err := fxr.backend.GetLowerPeopleList(req.UnionId, req.Offset, req.Num)
		if err != nil {
			logrus.Errorf("Error get fx lower people list: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx lower people list: %v", err)
		} else {
			var listInfo FxLowerPeopleList
			listInfo.Count = count
			listInfo.List = list
			rsp.Data = listInfo
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxAccountRank(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}
	
	req := &getFxAccountRankReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	
	rsp := &FxResponse{Code: RspCodeOK}
	
	list, err := fxr.backend.GetFxAccountRank(req.Offset, req.Num)
	if err != nil {
		logrus.Errorf("Error get fx account rank list: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx account rank list: %v", err)
	} else {
		rsp.Data = list
	}
	
	return utils.WriteJSON(w, http.StatusOK, rsp)
}


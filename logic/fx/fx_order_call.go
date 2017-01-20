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

func (fxr *FXRouter) createFxOrder(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &createFxOrderReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}
	//rsp.Msg = "暂时无法提交"
	//rsp.Code = RspCodeErr
	//return utils.WriteJSON(w, http.StatusOK, rsp)

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	err := fxr.backend.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("Error change union_id to account_id error: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error change union_id to account_id error: %v", err)
	} else {
		fxOrder := &models.FxOrder{
			AccountId:   fxAccount.ID,
			UnionId:     req.UnionId,
			OrderId:     req.OrderId,
			GoodsId:     req.GoodsId,
			OrderName:   req.OrderName,
			Price:       req.Price,
			ReturnMoney: req.ReturnMoney,
			Status:      req.Status,
		}
		if err := fxr.backend.CreateFxOrder(fxOrder); err != nil {
			logrus.Errorf("Error create fx order: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error create fx order: %v", err)
		}
	}
	rsp.Data = req.UnionId

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxOrderList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxOrderListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	err := fxr.backend.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("Error change union_id to account_id error: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error change union_id to account_id error: %v", err)
	} else {
		type FxOrderList struct {
			Count int64            `json:"count"`
			List  []models.FxOrder `json:"list"`
		}
		count, err := fxr.backend.GetFxOrderListCountById(fxAccount.ID)
		if err != nil {
			logrus.Errorf("Error get fx order list count: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx order list count: %v", err)
		} else {
			list, err := fxr.backend.GetFxOrderListById(fxAccount.ID, req.Offset, req.Num, req.Status)
			if err != nil {
				logrus.Errorf("Error get fx order list: %v", err)
				rsp.Code = RspCodeErr
				rsp.Msg = fmt.Sprintf("Error get fx order list: %v", err)
			} else {
				var listInfo FxOrderList
				listInfo.Count = count
				listInfo.List = list
				rsp.Data = listInfo
			}
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxOrderSettlementRecordList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxOrderSettlementRecordListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	err := fxr.backend.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("Error change union_id to account_id error: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error change union_id to account_id error: %v", err)
	} else {
		type FxOrderSettlementRecordList struct {
			Count int64                            `json:"count"`
			List  []models.FxOrderSettlementRecord `json:"list"`
		}
		count, err := fxr.backend.GetFxOrderSettlementRecordListCountById(fxAccount.ID)
		if err != nil {
			logrus.Errorf("Error get fx order settlement record list count: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx order settlement record list count: %v", err)
		} else {
			list, err := fxr.backend.GetFxOrderSettlementRecordListById(fxAccount.ID, req.Offset, req.Num)
			if err != nil {
				logrus.Errorf("Error get fx order settlement record list: %v", err)
				rsp.Code = RspCodeErr
				rsp.Msg = fmt.Sprintf("Error get fx order settlement record list: %v", err)
			} else {
				var listInfo FxOrderSettlementRecordList
				listInfo.Count = count
				listInfo.List = list
				rsp.Data = listInfo
			}
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxOrderWaitSettlementSum(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxOrderWaitSettlementSumReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	err := fxr.backend.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("Error change union_id to account_id error: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error change union_id to account_id error: %v", err)
	} else {
		sum, err := fxr.backend.GetFxOrderWaitSettlementRecordSum(fxAccount.ID)
		if err != nil {
			logrus.Errorf("Error get fx order wait settlement sum: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx order wait settlement sum: %v", err)
		} else {
			rsp.Data = sum
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxOrderWaitSettlementRecordList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxOrderWaitSettlementRecordListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	err := fxr.backend.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("Error change union_id to account_id error: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error change union_id to account_id error: %v", err)
	} else {
		type FxOrderWaitSettlementRecordList struct {
			Count int64                                `json:"count"`
			List  []models.FxOrderWaitSettlementRecord `json:"list"`
		}
		count, err := fxr.backend.GetFxOrderWaitSettlementRecordListCountById(fxAccount.ID)
		if err != nil {
			logrus.Errorf("Error get fx order wait settlement record list count: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx order wait settlement record list count: %v", err)
		} else {
			list, err := fxr.backend.GetFxOrderWaitSettlementRecordListById(fxAccount.ID, req.Offset, req.Num)
			if err != nil {
				logrus.Errorf("Error get fx order wait settlement record list: %v", err)
				rsp.Code = RspCodeErr
				rsp.Msg = fmt.Sprintf("Error get fx order wait settlement record list: %v", err)
			} else {
				var listInfo FxOrderWaitSettlementRecordList
				listInfo.Count = count
				listInfo.List = list
				rsp.Data = listInfo
			}
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

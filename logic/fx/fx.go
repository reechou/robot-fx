package fx

import (
	"github.com/reechou/robot-fx/config"
	"github.com/reechou/robot-fx/router"
)

type FXRouter struct {
	cfg     *config.Config
	backend Backend
	routes  []router.Route
}

func NewRouter(b Backend, cfg *config.Config) router.Router {
	r := &FXRouter{
		cfg:     cfg,
		backend: b,
	}
	r.initRoutes()
	return r
}

func (fxr *FXRouter) Routes() []router.Route {
	return fxr.routes
}

func (fxr *FXRouter) initRoutes() {
	fxr.routes = []router.Route{
		router.NewPostRoute("/create_robot_alimamam", fxr.createRobotAlimama),
		router.NewPostRoute("/robot_call", fxr.robotCall),

		// about fx account
		router.NewPostRoute("/fx/create_fx_account", fxr.createFxAccount),
		router.NewPostRoute("/fx/create_fx_salesman", fxr.createFxSalesman),
		router.NewPostRoute("/fx/update_fx_baseinfo", fxr.updateFxBaseInfo),
		router.NewPostRoute("/fx/update_fx_status", fxr.updateFxStatus),
		router.NewPostRoute("/fx/fx_sign", fxr.updateFxSignTime),
		router.NewPostRoute("/fx/get_fx_accout", fxr.getFxAccount),
		router.NewPostRoute("/fx/get_fx_lower_people_list", fxr.getFxLowerPeopleList),
		router.NewPostRoute("/fx/get_fx_account_rank", fxr.getFxAccountRank),
		// about fx account history
		router.NewPostRoute("/fx/get_fx_history", fxr.getFxAccountHistoryList),
		router.NewPostRoute("/fx/get_fx_history_by_type", fxr.getFxAccountHistoryListByType),
		// about fx order
		router.NewPostRoute("/fx/create_fx_order", fxr.createFxOrder),
		router.NewPostRoute("/fx/get_fx_order_list", fxr.getFxOrderList),
		router.NewPostRoute("/fx/get_fx_order_wait_sr_sum", fxr.getFxOrderWaitSettlementSum),
		router.NewPostRoute("/fx/get_fx_order_sr_list", fxr.getFxOrderSettlementRecordList),
		router.NewPostRoute("/fx/get_fx_order_wait_sr_list", fxr.getFxOrderWaitSettlementRecordList),
		// about withdrawal
		router.NewPostRoute("/fx/create_fx_withdrawal", fxr.createFxWithdrawalRecord),
		router.NewPostRoute("/fx/get_fx_withdrawal_record", fxr.getFxWithdrawalRecordList),
		router.NewPostRoute("/fx/get_fx_withdrawal_sum", fxr.getFxWithdrawalRecordSum),
		router.NewPostRoute("/fx/get_fx_withdrawal_error_list", fxr.getFxWithdrawalRecordErrorList),
		router.NewPostRoute("/fx/get_fx_withdrawal_error_list_from_name", fxr.getFxWithdrawalRecordErrorListFromName),
		router.NewPostRoute("/fx/get_fx_withdrawal_all", fxr.getFxWithdrawalRecordAll),
		router.NewPostRoute("/fx/confirm_withdrawal", fxr.confirmWithdrawalRecord),
		
		router.NewOptionsRoute("/fx/get_fx_withdrawal_all", fxr.optionCall),
		router.NewOptionsRoute("/fx/confirm_withdrawal", fxr.optionCall),
	}
}

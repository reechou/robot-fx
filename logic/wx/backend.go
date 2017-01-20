package wx

import (
	"github.com/reechou/robot-fx/logic/models"
)

type wxCallBackend interface {
	WXCheck(info *models.WXCallCheck) error
	WXHandleReq(req *models.WXCallRequest) (*models.WXCallResponse, error)
}

type Backend interface {
	wxCallBackend
}

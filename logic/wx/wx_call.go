package wx

import (
	"encoding/xml"
	"net/http"

	"github.com/reechou/robot-fx/logic/models"
	"github.com/reechou/robot-fx/utils"
	"golang.org/x/net/context"
)

func (wxr *WXRouter) wxCallGet(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}
	call := &models.WXCallCheck{
		Signature: r.Form.Get("signature"),
		Timestamp: r.Form.Get("timestamp"),
		Nonce:     r.Form.Get("nonce"),
		Echostr:   r.Form.Get("echostr"),
	}
	err := wxr.backend.WXCheck(call)
	if err != nil {
		w.Write([]byte("not success."))
		return err
	}
	w.Write([]byte(call.Echostr))

	return nil
}

func (wxr *WXRouter) wxCallPost(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &models.WXCallRequest{}
	if err := xml.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp, err := wxr.backend.WXHandleReq(req)
	if err != nil {
		return nil
	}

	return utils.WriteXML(w, http.StatusOK, rsp)
}

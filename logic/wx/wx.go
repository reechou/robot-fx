package wx

import (
	"github.com/reechou/robot-fx/router"
)

type WXRouter struct {
	backend Backend
	routes  []router.Route
}

func NewRouter(b Backend) router.Router {
	r := &WXRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

func (wxr *WXRouter) Routes() []router.Route {
	return wxr.routes
}

func (wxr *WXRouter) initRoutes() {
	wxr.routes = []router.Route{
		router.NewGetRoute("/wx/call", wxr.wxCallGet),
		router.NewPostRoute("/wx/call", wxr.wxCallPost),
	}
}

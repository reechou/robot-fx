package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/reechou/robot-fx/middleware"
	"github.com/reechou/robot-fx/utils"
)

// handleWithGlobalMiddlwares wraps the handler function for a request with
// the server's global middlewares. The order of the middlewares is backwards,
// meaning that the first in the list will be evaluated last.
func (s *Server) handleWithGlobalMiddlewares(handler utils.APIFunc) utils.APIFunc {
	next := handler

	if s.cfg.EnableCors {
		handleCORS := middleware.NewCORSMiddleware(s.cfg.CorsHeaders)
		next = handleCORS(next)
	}

	// Only want this on debug level
	if s.cfg.Logging && logrus.GetLevel() == logrus.DebugLevel {
		next = middleware.DebugRequestMiddleware(next)
	}

	return next
}

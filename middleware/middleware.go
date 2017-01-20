package middleware

import "github.com/reechou/robot-fx/utils"

// Middleware is an adapter to allow the use of ordinary functions as server API filters.
// Any function that has the appropriate signature can be register as a middleware.
type Middleware func(handler utils.APIFunc) utils.APIFunc

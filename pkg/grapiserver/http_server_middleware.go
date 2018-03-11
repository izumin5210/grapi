package grapiserver

import "net/http"

// HTTPServerMiddleware is an interface of http server middleware
type HTTPServerMiddleware func(http.Handler) http.Handler

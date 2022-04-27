package infrastructure

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/context"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/interfaces/handlers"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/interfaces/loggers"
	"gopkg.in/gorilla/mux.v1"
)

// Route stands for an http endpoint description
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler handlers.Handler
}

type routeGroups struct {
	Prefix string
	Groups []Route
}

// WrapperFunc defines a type for functions that wrap an http.HandlerFunc
// to modify its behaviour
type WrapperFunc func(pattern string, handler http.HandlerFunc) http.HandlerFunc

// Routes is an array of routes with a common prefix
type Routes []routeGroups

// RouterMaker gathers route and wrapper information to build a router
type RouterMaker struct {
	Logger        loggers.Logger
	WrapperFuncs  []WrapperFunc
	WithProfiling bool
	Routes        Routes
}

// NewRouter setups a Router based on the provided routes
func (maker *RouterMaker) NewRouter() http.Handler {
	router := mux.NewRouter()
	for _, routeGroup := range maker.Routes {
		subRouter := router.PathPrefix(routeGroup.Prefix).Subrouter()
		for _, route := range routeGroup.Groups {
			hLogger := loggers.MakeJSONHandlerLogger(maker.Logger)
			handler := handlers.MakeJSONHandlerFunc(route.Handler, hLogger)
			for _, wrapFunc := range maker.WrapperFuncs {
				handler = wrapFunc(route.Pattern, handler)
			}
			subRouter.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
	if maker.WithProfiling {
		router.HandleFunc("/debug/pprof/", pprof.Index)
		router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		router.HandleFunc("/debug/pprof/profile", pprof.Profile)
		router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		router.HandleFunc("/debug/pprof/trace", pprof.Trace)

		router.Handle("/debug/pprof/block", pprof.Handler("block"))
		router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
		router.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
		router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	}
	return context.ClearHandler(router)
}

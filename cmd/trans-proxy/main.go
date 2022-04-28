package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/infrastructure"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/interfaces/handlers"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/interfaces/loggers"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/interfaces/repository/services"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/usecases"
)

var shutdownSequence = infrastructure.NewShutdownSequence()

func main() { // nolint funlen
	var conf infrastructure.Config
	shutdownSequence.Listen()
	infrastructure.LoadFromEnv(&conf)
	if jconf, err := json.MarshalIndent(conf, "", "    "); err == nil {
		fmt.Printf("Config: \n%s\n", jconf)
	}

	fmt.Printf("Setting up Prometheus\n")
	prometheus := infrastructure.MakePrometheusExporter(
		conf.PrometheusConf.Port,
		conf.PrometheusConf.Enabled,
	)

	fmt.Printf("Setting up logger\n")
	logger, err := infrastructure.MakeYapoLogger(&conf.LoggerConf,
		prometheus.NewEventsCollector(
			"trans-proxy_service_events_total",
			"events tracker counter for trans-proxy service",
		),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	logger.Info("Initializing resources")
	useBrowserCache := infrastructure.InBrowserCache{
		MaxAge:  conf.InBrowserCacheConf.MaxAge,
		Etag:    conf.InBrowserCacheConf.Etag,
		Enabled: conf.InBrowserCacheConf.Enabled,
	}
	// HealthHandler
	var healthHandler handlers.HealthHandler

	// transHandler
	transFactory := infrastructure.NewTextProtocolTransFactory(conf.Trans, logger)
	transRepository := services.NewTransRepo(transFactory)
	transLogger := loggers.MakeTransInteractorLogger(logger)
	transInteractor := usecases.TransInteractor{
		Repository: transRepository,
		Logger:     transLogger,
	}
	transHandler := handlers.TransHandler{
		Interactor: transInteractor,
		TokenValidationInteractor: &usecases.ValidateToken{
			SecretToken: conf.Runtime.APIKey,
		},
	}
	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger:         logger,
		Cors:           conf.CorsConf,
		InBrowserCache: useBrowserCache,
		WrapperFuncs: []infrastructure.WrapperFunc{
			prometheus.TrackHandlerFunc,
		},
		WithProfiling: conf.Runtime.Profiling,
		Routes: infrastructure.Routes{
			{
				// This is the base path, all routes will start with this prefix
				Prefix: "/api/v{version:[1-9][0-9]*}",
				Groups: []infrastructure.Route{
					{
						Name:         "Check service health",
						Method:       "GET",
						Pattern:      "/healthcheck",
						Handler:      &healthHandler,
						RequestCache: "10s",
					},
					{
						Name:    "Execute a trans request",
						Method:  "POST",
						Pattern: "/execute/{command}",
						Handler: &transHandler,
					},
				},
			},
		},
	}
	server := infrastructure.NewHTTPServer(
		conf.Runtime.Address(),
		maker.NewRouter(),
		logger,
	)
	shutdownSequence.Push(server)
	logger.Info("Starting request serving")
	go server.ListenAndServe()
	shutdownSequence.Wait()
	logger.Info("Server exited normally")
}

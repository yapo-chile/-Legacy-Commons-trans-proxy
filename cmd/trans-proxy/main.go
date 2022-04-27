package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/infrastructure"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/interfaces/handlers"
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
	}

	trans-proxyHandler := handlers.TransHandler{
		Interactor: trans-proxyInteractor,
	}
	// Setting up router
	maker := infrastructure.RouterMaker{
		Logger: logger,
		WrapperFuncs: []infrastructure.WrapperFunc{
			prometheus.TrackHandlerFunc,
		},
		WithProfiling: conf.ServiceConf.Profiling,
		Routes: infrastructure.Routes{
			{
						Pattern: "/healthcheck",
						Handler: &healthHandler,
					},
					{
						Name:    "Execute a trans-proxy request",
						Method:  "POST",
						Pattern: "/execute/{command}",
						Handler: &trans-proxyHandler,
					},
				},
			},
		},
	}
	server := infrastructure.NewHTTPServer(
		fmt.Sprintf("%s:%d", conf.Runtime.Host, conf.Runtime.Port),
		maker.NewRouter(),
		logger,
	)
	shutdownSequence.Push(server)
	go server.ListenAndServe()
	shutdownSequence.Wait()

	logger.Info("Starting request serving")
	logger.Crit("%s\n", http.ListenAndServe(conf.ServiceConf.Host, maker.NewRouter()))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/artem-shestakov/autofaq_exporter/internal/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func main() {
	port := flag.String("listen-address", ":9901", "Address on which to bind and expose metrics")
	autofaqUrl := flag.String("autofaq", "", "URL address of AutoFAQ server")
	logLevel := flag.String("log-level", "info", "Set log level output. Where it is one of 'debug','info','warn','error'")
	flag.Parse()

	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestamp)
	switch *logLevel {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "info":
		logger = level.NewFilter(logger, level.AllowInfo())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		level.Warn(logger).Log("msg", fmt.Sprintf("Unknown log level '%s'", *logLevel))
		level.Warn(logger).Log("msg", "Log level set to default value 'info'")
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	autoFAQCollector, err := collector.NewAutoFAQCollector(*autofaqUrl, logger)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
		log.Panic(err.Error())
	}
	level.Debug(logger).Log("msg", "AutoFAQ collector is created")

	prometheus.MustRegister(autoFAQCollector)
	http.Handle("/metrics", promhttp.Handler())

	level.Info(logger).Log("msg", fmt.Sprintf("autofaq_exporter starting on %s", *port))
	http.ListenAndServe(*port, nil)
}

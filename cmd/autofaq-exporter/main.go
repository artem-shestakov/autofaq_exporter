package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/artem-shestakov/autofaq_exporter/internal/collector"
	"github.com/artem-shestakov/autofaq_exporter/internal/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func main() {
	// Init logger
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestamp)

	// Flags
	configPath := pflag.StringP("config", "c", "", "Path to config file")
	pflag.Parse()

	// Init config
	_, err := config.InitConfig(configPath, logger)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
	}
	// Set vars from config or Env variables
	listenAddress := viper.GetString("listen-address")
	if listenAddress == "" {
		listenAddress = ":9901"
	}
	autofaqUrl := viper.GetString("autofaq-url")
	if autofaqUrl == "" {
		level.Error(logger).Log("msg", "AutoFAQ URL is empty, please set URL")
		log.Fatal()
	}

	// Set log level
	switch viper.GetString("log-level") {
	case "debug":
		logger = level.NewFilter(logger, level.AllowDebug())
	case "info":
		logger = level.NewFilter(logger, level.AllowInfo())
	case "warn":
		logger = level.NewFilter(logger, level.AllowWarn())
	case "error":
		logger = level.NewFilter(logger, level.AllowError())
	default:
		level.Warn(logger).Log("msg", fmt.Sprintf("Unknown log level '%s'", viper.GetString("severity")))
		level.Warn(logger).Log("msg", "Log level set to default value 'info'")
		logger = level.NewFilter(logger, level.AllowInfo())
	}

	// Create metric collector
	autoFAQCollector, err := collector.NewAutoFAQCollector(autofaqUrl, logger)
	if err != nil {
		level.Error(logger).Log("msg", err.Error())
		log.Panic(err.Error())
	}
	level.Debug(logger).Log("msg", "AutoFAQ collector is created")
	prometheus.MustRegister(autoFAQCollector)

	// Start exporter
	http.Handle("/metrics", promhttp.Handler())
	level.Info(logger).Log("msg", fmt.Sprintf("autofaq_exporter starting on %s", listenAddress))
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		level.Error(logger).Log("msg", err.Error())
	}
}

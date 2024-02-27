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
	port := flag.String("listen-address", ":9901", "address on which to bind and expose metrics")
	autofaqUrl := flag.String("autofaq", "", "address of AutoFAQ server")

	flag.Parse()
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestamp)
	logger = level.NewFilter(logger, level.AllowInfo())

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

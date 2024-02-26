package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/artem-shestakov/autofaq_exporter/internal/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := flag.String("listen-address", ":9901", "address on which to bind and expose metrics")
	autofaqUrl := flag.String("autofaq", "", "address of AutoFAQ server")

	flag.Parse()

	autoFAQCollector, err := collector.NewAutoFAQCollector(*autofaqUrl)
	if err != nil {
		log.Panic(err.Error())
	}
	prometheus.MustRegister(autoFAQCollector)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(*port, nil)
}

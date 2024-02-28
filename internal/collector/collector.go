package collector

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var initFuncs = make(map[string]func() (Collector, error))

func registerCollector(collector string, initFunc func() (Collector, error)) {
	initFuncs[collector] = initFunc
}

type Collector interface {
	Update(autofaq string, ch chan<- prometheus.Metric) error
}

type AutoFAQCollector struct {
	AutoFAQURL         string
	Collectors         map[string]Collector
	Logger             kitlog.Logger
	scrapeDurationDesc *prometheus.Desc
	scrapeSuccessDesc  *prometheus.Desc
}

func (a AutoFAQCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- a.scrapeDurationDesc
	ch <- a.scrapeSuccessDesc
}

func (a AutoFAQCollector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(a.Logger).Log("msg", "Start collect metrics")
	wg := sync.WaitGroup{}
	wg.Add(len(a.Collectors))
	for name, c := range a.Collectors {
		go func(name string, c Collector) {
			level.Debug(a.Logger).Log("msg", fmt.Sprintf("Start collect metrics of '%s' collector", name))
			a.execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
	level.Debug(a.Logger).Log("msg", "Finish collect metrics")
}

func (a AutoFAQCollector) execute(name string, c Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Update(a.AutoFAQURL, ch)
	duration := time.Since(begin)
	var success float64
	if err != nil {
		success = 0
	} else {
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(a.scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name, a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.scrapeSuccessDesc, prometheus.GaugeValue, success, name, a.AutoFAQURL)
}

func NewAutoFAQCollector(autofaq string, logger kitlog.Logger) (*AutoFAQCollector, error) {
	collectors := make(map[string]Collector)
	for name, initFunc := range initFuncs {
		collector, err := initFunc()
		if err != nil {
			return nil, err
		}
		collectors[name] = collector
	}
	return &AutoFAQCollector{
		AutoFAQURL: autofaq,
		Collectors: collectors,
		Logger:     logger,
		scrapeDurationDesc: prometheus.NewDesc("collector_duration_seconds",
			"autofaq_exporter: Duration of a collector scrape",
			[]string{"collector", "site"}, nil),
		scrapeSuccessDesc: prometheus.NewDesc("collector_success",
			"autofaq_exporter: Whether a collector succeeded",
			[]string{"collector", "site"}, nil),
	}, nil
}

package collector

import (
	"github.com/artem-shestakov/autofaq_exporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var initFuncs = make(map[string]func() (Collector, error))

type Collector interface {
	Update(autofaq string, logger kitlog.Logger, ch chan<- prometheus.Metric) error
}

// Collector
type AutoFAQCollector struct {
	AutoFAQURL            string
	Collectors            map[string]Collector
	Logger                kitlog.Logger
	Services              []config.Service
	scrapeDurationDesc    *prometheus.Desc
	scrapeSuccessDesc     *prometheus.Desc
	upTime                *prometheus.Desc
	dbUp                  *prometheus.Desc
	totalConnections      *prometheus.Desc
	activeConnections     *prometheus.Desc
	idleConnections       *prometheus.Desc
	runtimeTotal          *prometheus.Desc
	runtimeFree           *prometheus.Desc
	runtimeUsed           *prometheus.Desc
	garbageCollectionTime *prometheus.Desc
	status                *prometheus.Desc
	widgetStatus          *prometheus.Desc
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (a AutoFAQCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- a.upTime
	ch <- a.dbUp
	ch <- a.totalConnections
	ch <- a.activeConnections
	ch <- a.idleConnections
	ch <- a.runtimeTotal
	ch <- a.runtimeFree
	ch <- a.runtimeUsed
	ch <- a.garbageCollectionTime
	ch <- a.status
	ch <- a.widgetStatus
	ch <- a.scrapeDurationDesc
	ch <- a.scrapeSuccessDesc
}

// Collect implements required collect function for all promehteus collectors
func (a AutoFAQCollector) Collect(ch chan<- prometheus.Metric) {
	level.Debug(a.Logger).Log("msg", "Start collect metrics")

	// Sys collect
	a.collectSysMetrics(ch)
	// Widget collect
	a.collectWidgetsMetrics(ch)
	level.Debug(a.Logger).Log("msg", "Finish collect metrics")
}

func NewAutoFAQCollector(autofaq string, logger kitlog.Logger, services []config.Service) (*AutoFAQCollector, error) {
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
		Services:   services,
		scrapeDurationDesc: prometheus.NewDesc("collector_duration_seconds",
			"autofaq_exporter: Duration of a collector scrape",
			[]string{"collector", "site"}, nil),
		scrapeSuccessDesc: prometheus.NewDesc("collector_success",
			"autofaq_exporter: Whether a collector succeeded",
			[]string{"collector", "site"}, nil),
		upTime: prometheus.NewDesc("autofaq_sys_uptime",
			"Show backend start time", []string{"site"}, nil),
		dbUp: prometheus.NewDesc("autofaq_sys_db_up",
			"Show if AutoFAQ database is up", []string{"site"}, nil),
		totalConnections: prometheus.NewDesc("autofaq_sys_total_conn",
			"Total connections to DB", []string{"site"}, nil),
		activeConnections: prometheus.NewDesc("autofaq_sys_active_conn",
			"Active connections to DB", []string{"site"}, nil),
		idleConnections: prometheus.NewDesc("autofaq_sys_idle_conn",
			"Idle connections to DB", []string{"site"}, nil),
		runtimeTotal: prometheus.NewDesc("autofaq_sys_runtime_total",
			"JVM runtime total memory", []string{"site"}, nil),
		runtimeFree: prometheus.NewDesc("autofaq_sys_runtime_free",
			"JVM tuntime free memory", []string{"site"}, nil),
		runtimeUsed: prometheus.NewDesc("autofaq_sys_runtime_used",
			"JVM tuntime used memory", []string{"site"}, nil),
		garbageCollectionTime: prometheus.NewDesc("autofaq_sys_garbage_collection_time",
			"JVM garbage collection time", []string{"site"}, nil),
		status: prometheus.NewDesc("autofaq_sys_status",
			"Show if AutoFAQ backend server is up", []string{"site"}, nil),
		widgetStatus: prometheus.NewDesc("autofaq_widget_status", "Status of widget",
			[]string{"site", "service_id", "widget_id"}, nil),
	}, nil
}

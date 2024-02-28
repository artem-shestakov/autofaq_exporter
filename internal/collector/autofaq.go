package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

type AutoFAQSysInfo struct {
	BuildInfo   BuildInfo   `json:"buildInfo"`
	DbInfo      DbInfo      `json:"dbInfo"`
	RuntimeInfo RuntimeInfo `json:"runtimeInfo"`
	Status      string      `json:"status"`
}

type BuildInfo struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	ScalaVersion    string `json:"scalaVersion"`
	SbtVersion      string `json:"sbtVersion"`
	BuildTimestamp  string `json:"buildTimestamp"`
	GitHash         string `json:"gitHash"`
	AutofaqUrl      string `json:"autofaqUrl"`
	AutofaqUrlCrud  string `json:"autofaqUrlCrud"`
	AutofaqUrlQuery string `json:"autofaqUrlQuery"`
	UpTime          int    `json:"upTime"`
	AuthType        string `json:"authType"`
}

type DbInfo struct {
	DbUp              string `json:"dbUp"`
	TotalConnections  int    `json:"totalConnections"`
	ActiveConnections int    `json:"activeConnections"`
	IdleConnections   int    `json:"idleConnections"`
}

type RuntimeInfo struct {
	Total                 int `json:"total"`
	Free                  int `json:"free"`
	Used                  int `json:"used"`
	GarbageCollectionTime int `json:"garbageCollectionTime"`
}

type AutoFAQSysInfoCollector struct {
	UpTime                *prometheus.Desc
	DbUp                  *prometheus.Desc
	TotalConnections      *prometheus.Desc
	ActiveConnections     *prometheus.Desc
	IdleConnections       *prometheus.Desc
	RuntimeTotal          *prometheus.Desc
	RuntimeFree           *prometheus.Desc
	RuntimeUsed           *prometheus.Desc
	GarbageCollectionTime *prometheus.Desc
	Status                *prometheus.Desc
}

// Parse sys info from AutoFAQ site
func (c *AutoFAQSysInfoCollector) getSysInfo(autofaqURL string) (*AutoFAQSysInfo, error) {
	var autoFAQSysInfo AutoFAQSysInfo
	resp, err := http.Get(fmt.Sprintf("%s/api/sysInfo", autofaqURL))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &autoFAQSysInfo)
	return &autoFAQSysInfo, nil
}

// Collect metrics and publish them
func (c *AutoFAQSysInfoCollector) Update(autoFAQSite string, logger kitlog.Logger, ch chan<- prometheus.Metric) error {
	var dbUp, status int
	level.Debug(logger).Log("msg", fmt.Sprintf("Parse data from '%s'", autoFAQSite))
	autoFAQSysInfo, err := c.getSysInfo(autoFAQSite)
	if autoFAQSysInfo.DbInfo.DbUp == "success" {
		dbUp = 1
	} else {
		dbUp = 0
	}
	if autoFAQSysInfo.Status == "success" {
		status = 1
	} else {
		status = 0
	}
	level.Debug(logger).Log("msg", "Publish metrics of 'autofaq_sys_info' collector")
	ch <- prometheus.MustNewConstMetric(c.UpTime, prometheus.GaugeValue, float64(autoFAQSysInfo.BuildInfo.UpTime), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.DbUp, prometheus.GaugeValue, float64(dbUp), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.TotalConnections, prometheus.GaugeValue, float64(autoFAQSysInfo.DbInfo.TotalConnections), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.ActiveConnections, prometheus.GaugeValue, float64(autoFAQSysInfo.DbInfo.ActiveConnections), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.IdleConnections, prometheus.GaugeValue, float64(autoFAQSysInfo.DbInfo.IdleConnections), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.RuntimeTotal, prometheus.GaugeValue, float64(autoFAQSysInfo.RuntimeInfo.Total), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.RuntimeUsed, prometheus.GaugeValue, float64(autoFAQSysInfo.RuntimeInfo.Used), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.RuntimeFree, prometheus.GaugeValue, float64(autoFAQSysInfo.RuntimeInfo.Free), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.GarbageCollectionTime, prometheus.GaugeValue, float64(autoFAQSysInfo.RuntimeInfo.GarbageCollectionTime), autoFAQSite)
	ch <- prometheus.MustNewConstMetric(c.Status, prometheus.GaugeValue, float64(status), autoFAQSite)
	return err
}

// Add collector as child collector
func init() {
	registerCollector("autofaq_sys_info", NewAutoFAQSysInfoCollector)
}

func NewAutoFAQSysInfoCollector() (Collector, error) {
	return &AutoFAQSysInfoCollector{
		UpTime: prometheus.NewDesc("autofaq_sys_uptime",
			"Show backend start time", []string{"site"}, nil),
		DbUp: prometheus.NewDesc("autofaq_sys_db_up",
			"Show if AutoFAQ database is up", []string{"site"}, nil),
		TotalConnections: prometheus.NewDesc("autofaq_sys_total_conn",
			"Total connections to DB", []string{"site"}, nil),
		ActiveConnections: prometheus.NewDesc("autofaq_sys_active_conn",
			"Active connections to DB", []string{"site"}, nil),
		IdleConnections: prometheus.NewDesc("autofaq_sys_idle_conn",
			"Idle connections to DB", []string{"site"}, nil),
		RuntimeTotal: prometheus.NewDesc("autofaq_sys_runtime_total",
			"JVM runtime total memory", []string{"site"}, nil),
		RuntimeFree: prometheus.NewDesc("autofaq_sys_runtime_free",
			"JVM tuntime free memory", []string{"site"}, nil),
		RuntimeUsed: prometheus.NewDesc("autofaq_sys_runtime_used",
			"JVM tuntime used memory", []string{"site"}, nil),
		GarbageCollectionTime: prometheus.NewDesc("autofaq_sys_garbage_collection_time",
			"JVM garbage collection time", []string{"site"}, nil),
		Status: prometheus.NewDesc("autofaq_sys_status",
			"Show if AutoFAQ backend server is up", []string{"site"}, nil),
	}, nil
}

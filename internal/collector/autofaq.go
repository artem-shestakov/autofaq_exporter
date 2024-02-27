package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func (c *AutoFAQSysInfoCollector) Update(autofaq string, ch chan<- prometheus.Metric) error {
	var dbUp, status int
	autoFAQSysInfo, err := c.getSysInfo(autofaq)
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
	ch <- prometheus.MustNewConstMetric(c.UpTime, prometheus.GaugeValue, float64(autoFAQSysInfo.BuildInfo.UpTime))
	ch <- prometheus.MustNewConstMetric(c.DbUp, prometheus.GaugeValue, float64(dbUp))
	ch <- prometheus.MustNewConstMetric(c.Status, prometheus.GaugeValue, float64(status))
	fmt.Println(autoFAQSysInfo)
	return err
}

func init() {
	registerCollector("autofaq_sys_info", NewAutoFAQSysInfoCollector)
}

func NewAutoFAQSysInfoCollector() (Collector, error) {
	return &AutoFAQSysInfoCollector{
		UpTime: prometheus.NewDesc("autofaq_sys_uptime",
			"Show backend uptime", nil, nil),
		DbUp: prometheus.NewDesc("autofaq_sys_db_up",
			"Show if AutoFAQ database is up", nil, nil),
		TotalConnections: prometheus.NewDesc("autofaq_sys_total_conn",
			"Total connections to DB", nil, nil),
		ActiveConnections: prometheus.NewDesc("autofaq_sys_active_conn",
			"Active connections to DB", nil, nil),
		IdleConnections: prometheus.NewDesc("autofaq_sys_idle_conn",
			"Idle connections to DB", nil, nil),
		RuntimeTotal: prometheus.NewDesc("autofaq_sys_runtime_total",
			"runtime total", nil, nil),
		RuntimeFree: prometheus.NewDesc("autofaq_sys_runtime_free",
			"tuntime_free", nil, nil),
		RuntimeUsed: prometheus.NewDesc("autofaq_sys_runtime_used",
			"tuntime_used", nil, nil),
		GarbageCollectionTime: prometheus.NewDesc("autofaq_sys_garbage_collection_time",
			"garbage_collection_time", nil, nil),
		Status: prometheus.NewDesc("autofaq_sys_status",
			"Show if AutoFAQ backend server is up", nil, nil),
	}, nil
}

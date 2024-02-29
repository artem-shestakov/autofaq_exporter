package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

// Parse sys info from AutoFAQ site
func (a AutoFAQCollector) getSysInfo() (*AutoFAQSysInfo, error) {
	var autoFAQSysInfo AutoFAQSysInfo
	resp, err := http.Get(fmt.Sprintf("%s/api/sysInfo", a.AutoFAQURL))
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
func (a AutoFAQCollector) collectSysMetrics(ch chan<- prometheus.Metric) error {
	var dbUp, status int
	var success = float64(1)
	begin := time.Now()
	level.Debug(a.Logger).Log("msg", fmt.Sprintf("Parse data from '%s'", a.AutoFAQURL))
	autoFAQSysInfo, err := a.getSysInfo()
	if err != nil {
		success = 0
	}
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
	duration := time.Since(begin)

	level.Debug(a.Logger).Log("msg", "Publish metrics of 'autofaq_sys_info' collector")
	ch <- prometheus.MustNewConstMetric(a.upTime, prometheus.GaugeValue, float64(autoFAQSysInfo.BuildInfo.UpTime), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.dbUp, prometheus.GaugeValue, float64(dbUp), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.totalConnections, prometheus.GaugeValue, float64(autoFAQSysInfo.DbInfo.TotalConnections), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.activeConnections, prometheus.GaugeValue, float64(autoFAQSysInfo.DbInfo.ActiveConnections), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.idleConnections, prometheus.GaugeValue, float64(autoFAQSysInfo.DbInfo.IdleConnections), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.runtimeTotal, prometheus.GaugeValue, float64(autoFAQSysInfo.RuntimeInfo.Total), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.runtimeUsed, prometheus.GaugeValue, float64(autoFAQSysInfo.RuntimeInfo.Used), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.runtimeFree, prometheus.GaugeValue, float64(autoFAQSysInfo.RuntimeInfo.Free), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.garbageCollectionTime, prometheus.GaugeValue, float64(autoFAQSysInfo.RuntimeInfo.GarbageCollectionTime), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.status, prometheus.GaugeValue, float64(status), a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), "autofaq_sys_collector", a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.scrapeSuccessDesc, prometheus.GaugeValue, success, "autofaq_sys_collector", a.AutoFAQURL)
	return err
}

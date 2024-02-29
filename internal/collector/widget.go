package collector

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (a AutoFAQCollector) collectWidgetsMetrics(ch chan<- prometheus.Metric) {
	var success = float64(1)
	begin := time.Now()
	for _, service := range a.Services {
		for _, widgetId := range service.Widgets {
			status, err := a.getWidgetStatus(service.Id, widgetId)
			if err != nil {
				success = 0
				ch <- prometheus.MustNewConstMetric(a.widgetStatus, prometheus.GaugeValue, 0, a.AutoFAQURL, service.Id, widgetId)
			}
			if status == 200 {
				ch <- prometheus.MustNewConstMetric(a.widgetStatus, prometheus.GaugeValue, 1, a.AutoFAQURL, service.Id, widgetId)
			} else {
				ch <- prometheus.MustNewConstMetric(a.widgetStatus, prometheus.GaugeValue, 0, a.AutoFAQURL, service.Id, widgetId)
			}
		}
	}
	duration := time.Since(begin)
	ch <- prometheus.MustNewConstMetric(a.scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), "autofaq_widget_collector", a.AutoFAQURL)
	ch <- prometheus.MustNewConstMetric(a.scrapeSuccessDesc, prometheus.GaugeValue, success, "autofaq_widget_collector", a.AutoFAQURL)
}

func (a AutoFAQCollector) getWidgetStatus(serviceId, widgetId string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/webhooks/widget/%s/%s/settings", a.AutoFAQURL, serviceId, widgetId))
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

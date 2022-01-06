package uip_cloudflare

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	countUpdate = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "updateip_cloudflare_update",
			Help: "Number of ip updated on Cloudflare provider.",
		},
	)

	// Status ...
	providerStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "updateip_cloudflare_status",
			Help: "cloudflare Providers status.",
		},
	)

	// Histo ...
	funcTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "updateip_cloudflare_func_time",
		Help: "Time taken to do ...",
	},
		[]string{"where"},
	)

	eventReceive = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "updateip_cloudflare_event_receive",
			Help: "Count of events received on Cloudflare.",
		},
	)
)

func (d *PCloudflare) RegistryMetrics() map[string][]interface{} {

	x := make(map[string][]interface{})
	x["counter"] = []interface{}{countUpdate}
	x["gauge"] = []interface{}{providerStatus}
	x["gaugeVec"] = []interface{}{funcTime}

	return x
}

// TimeTrackS ...
func timeTrackS(start time.Time, name string) {
	elapsed := time.Since(start)
	funcTime.WithLabelValues(name).Observe(elapsed.Seconds())
}

package uip_ovh

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	countUpdate = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "updateip_ovh_update",
			Help: "Number of ip updated on OVH provider.",
		},
	)

	// Status ...
	providerStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "updateip_ovh_status",
			Help: "cloudflare Providers status.",
		},
	)

	// Histo ...
	funcTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "updateip_ovh_func_time",
		Help: "Time taken to do ...",
	},
		[]string{"where"},
	)

	eventReceive = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "updateip_ovh_event_receive",
			Help: "Count of events received on OVH Provider.",
		},
	)
)

func (d *Povh) RegistryMetrics() map[string][]interface{} {

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

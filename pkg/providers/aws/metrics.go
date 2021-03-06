package uip_aws

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	countUpdate = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "updateip_aws_update",
			Help: "Number of ip updated on AWS provider.",
		},
	)

	// Status ...
	providerStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "updateip_aws_status",
			Help: "AWS Providers status.",
		},
	)

	// Histo ...
	funcTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "updateip_aws_func_time",
		Help: "Time taken to do ...",
	},
		[]string{"where"},
	)

	eventReceive = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "updateip_aws_event_receive",
			Help: "Count of events received on AWS Provider.",
		},
	)
)

func (d *Paws) RegistryMetrics() map[string][]interface{} {

	x := make(map[string][]interface{})
	x["counter"] = []interface{}{countUpdate, eventReceive}
	x["gauge"] = []interface{}{providerStatus}
	x["gaugeVec"] = []interface{}{funcTime}

	return x
}

// TimeTrackS ...
func timeTrackS(start time.Time, name string) {
	elapsed := time.Since(start)
	funcTime.WithLabelValues(name).Observe(elapsed.Seconds())
}

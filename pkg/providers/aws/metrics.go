package uip_aws

import "github.com/prometheus/client_golang/prometheus"

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
)

func (d *Paws) RegistryMetrics() map[string][]interface{} {

	x := make(map[string][]interface{})
	x["counter"] = []interface{}{countUpdate}
	x["gauge"] = []interface{}{providerStatus}

	return x
}

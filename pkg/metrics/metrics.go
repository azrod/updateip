package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/azrod/updateip/pkg/config"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Metrics struct {
	Counters *map[string]prometheus.Counter
	Gauges   *map[string]*prometheus.GaugeVec
	registry *prometheus.Registry
	cfg      config.CFGMetrics
}

// NewMetrics returns a new Metrics struct
func Init(cfg config.CFGMetrics) *Metrics {

	return &Metrics{
		cfg:      cfg,
		registry: prometheus.NewRegistry(),
		Gauges:   &map[string]*prometheus.GaugeVec{},
		Counters: &map[string]prometheus.Counter{},
	}
}

// RUn
func (m *Metrics) Run() {
	m.registerMetrics()
	go m.hTTPServer()
}

func (m *Metrics) registerMetrics() {
	m.registry.MustRegister(collectors.NewBuildInfoCollector())
}

func (m *Metrics) RegisterPkg(rg map[string][]interface{}) {

	for _, v := range rg["gauge"] {
		m.registry.MustRegister(v.(prometheus.Gauge))
	}

	for _, v := range rg["counter"] {
		m.registry.MustRegister(v.(prometheus.Counter))
	}

	for _, v := range rg["gaugeVec"] {
		m.registry.MustRegister(v.(*prometheus.HistogramVec))
	}

	m.registry.MustRegister(collectors.NewGoCollector())
}

func (m *Metrics) hTTPServer() {

	if m.cfg.Path == "" {
		m.cfg.Path = "/metrics"
	}

	if m.cfg.Port == 0 {
		m.cfg.Port = 8080
	}

	if m.cfg.Host == "" {
		m.cfg.Host = "0.0.0.0"
	}

	r := mux.NewRouter()

	srv := &http.Server{
		Addr: m.cfg.Host + ":" + strconv.Itoa(m.cfg.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	if m.cfg.Logging {
		log.Warn().Msg("Logging metrics enabled. This is not recommended for production.")
		r.Use(loggingMiddleware)
	}

	r.Handle(m.cfg.Path, promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})).Methods("GET")
	log.Debug().Msgf("Metrics server listening on %s:%d", m.cfg.Host, m.cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("Metrics server error")
	}

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Info().Str("Method", r.Method).Str("URL", r.URL.String()).Msg("Request")
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

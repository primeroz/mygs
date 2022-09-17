package main

import (
	"flag"

	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

func main() {
	// =====================
	// Get OS parameter
	// =====================
	var bind string
	flag.StringVar(&bind, "bind", "0.0.0.0:9104", "bind")

	flag.Parse()

	// ========================
	// Regist handler
	// ========================
	prometheus.Register(version.NewCollector("query_exporter"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
      <html>
      <head><title>PortScan Exporter Metrics</title></head>
      <body>
      <h1>links</h1>
      <p><a href='/metrics'>Metrics</a></p>
      </body>
      </html>
      `))
	})

	// Register metrics http handler
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	// register http healthcheck
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	// start server
	log.Infof("Starting http server - %s", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
	}
}

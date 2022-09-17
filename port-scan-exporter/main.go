package main

import (
	"flag"

	"net/http"

	"github.com/jasonlvhit/gocron"
	"github.com/kelseyhightower/memkv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var store memkv.Store

// ========================
// Collect Pods and Scan
// ========================
// TODO Need lock to prevent multiple runs if it takes longer then interval ?
func collectAndScan() {
	collectPods()
	scanPods()
}

func scheduleCollectAndScan(interval uint64) {
	gocron.Every(interval).Minute().Do(collectAndScan)
	<-gocron.Start()
}

func main() {

	// Create Store
	store = memkv.New()

	// =====================
	// Get OS parameter
	// =====================
	var bind string
	var interval uint64
	var debuglog bool

	flag.StringVar(&bind, "bind", "0.0.0.0:9104", "bind address")
	flag.BoolVar(&debuglog, "debug", false, "enable debug log")
	flag.Uint64Var(&interval, "collect-interval-min", 5, "interval in minutes to perform Collect of Pods and Port Scan")

	flag.Parse()

	if debuglog {
		log.SetLevel(log.DebugLevel)
	}

	// ========================
	// HTTP handlers
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

	// ========================
	// start scheduler
	// ========================
	go collectAndScan() // First Scan at startup
	go scheduleCollectAndScan(interval)

	// ========================
	// start server
	// ========================
	log.Infof("Starting http server - %s", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
	}
}

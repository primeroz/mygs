package main

import (
	"flag"
	"fmt"
	"time"

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
func collectAndScan(min int, max int) {
	firstStart := time.Now()

	start := time.Now()
	collectPods()
	timeTrack(start, "Collecting Pods")

	start = time.Now()
	scanPods(min, max)
	time.Sleep(5 * time.Second) // Ugly Hack to allow go subroutines to finish
	timeTrack(start, "Scanning Pods")

	store.Set("/timings/duration", fmt.Sprintf("%f", time.Since(firstStart).Seconds()))
	store.Set("/timings/last", fmt.Sprintf("%d", time.Now().Unix()))

	// // DEBUG Code to show content of Store
	// for _, pod := range store.ListDir("/ports") {
	// 	name, _ := store.GetValue(fmt.Sprintf("/pods/%s/name", pod))
	// 	namespace, _ := store.GetValue(fmt.Sprintf("/pods/%s/namespace", pod))
	// 	log.Debugf("Open ports for pod %s in namespace %s", name, namespace)
	// 	for _, port := range store.List(fmt.Sprintf("/ports/%s/", pod)) {
	// 		log.Debugf("  Port %s is open", port)
	// 	}
	// }
}

func scheduleCollectAndScan(interval uint64, min int, max int) {
	gocron.Every(interval).Minute().Do(collectAndScan, min, max)
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
	var portMin int
	var portMax int
	var debuglog bool

	flag.StringVar(&bind, "bind", "0.0.0.0:9104", "bind address")
	flag.BoolVar(&debuglog, "debug", false, "enable debug log")
	flag.Uint64Var(&interval, "collect-interval-min", 5, "interval in minutes to perform Collect of Pods and Port Scan")
	flag.IntVar(&portMin, "port-min", 1, "Min port to scan for")
	flag.IntVar(&portMax, "port-max", 10000, "Max port to scan for")

	flag.Parse()

	if debuglog {
		log.SetLevel(log.DebugLevel)
	}

	// ========================
	// HTTP handlers
	// ========================
	prometheus.Register(version.NewCollector("query_exporter"))

	// Register the Port Scanner prometheus exporter
	portScanRegister()

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
	go collectAndScan(portMin, portMax) // First Scan at startup
	go scheduleCollectAndScan(interval, portMin, portMax)

	// ========================
	// start server
	// ========================
	log.Infof("Starting http server - %s", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
	}
}

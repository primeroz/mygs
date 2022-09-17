package main

import (
	"fmt"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	//								     log "github.com/sirupsen/logrus"
)

// Define a struct for you collector that contains pointers
// to prometheus descriptors for each metric you wish to expose.
// Note you can also include fields of other types if they provide utility
// but we just won't be exposing them as metrics.
type portScanCollector struct {
	openPort           *prometheus.Desc
	scanDuration       *prometheus.Desc
	lastSuccessfulScan *prometheus.Desc
	scannedPodsTotal   *prometheus.Desc
	openPodsTotal      *prometheus.Desc
}

// You must create a constructor for you collector that
// initializes every descriptor and returns a pointer to the collector
func newPortScanCollector() *portScanCollector {
	return &portScanCollector{
		openPort: prometheus.NewDesc(prometheus.BuildFQName("port_scanner", "", "open_port"),
			"Open port for a pod",
			[]string{"pod", "namespace", "port"}, nil,
		),
		scanDuration: prometheus.NewDesc(prometheus.BuildFQName("port_scanner", "", "scan_duration_seconds"),
			"Duration of the Collect and PortScan phases",
			[]string{}, nil,
		),
		lastSuccessfulScan: prometheus.NewDesc(prometheus.BuildFQName("port_scanner", "", "last_successful_scan_epoch"),
			"Epoch of last successful scan",
			[]string{}, nil,
		),
		scannedPodsTotal: prometheus.NewDesc(prometheus.BuildFQName("port_scanner", "", "scanned_pods_total"),
			"Total number of scanned pods in the last scan",
			[]string{}, nil,
		),
		openPodsTotal: prometheus.NewDesc(prometheus.BuildFQName("port_scanner", "", "open_pods_total"),
			"Total number of pods with open ports in the last scan",
			[]string{}, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *portScanCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.openPort
	ch <- collector.scanDuration
}

// Collect implements required collect function for all promehteus collectors
func (collector *portScanCollector) Collect(ch chan<- prometheus.Metric) {
	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.

	// Create metrics for open ports - GAUGE value since it just represent an absolute value 1 for summing purposes
	for _, pod := range store.ListDir("/ports") {
		name, _ := store.GetValue(fmt.Sprintf("/pods/%s/name", pod))
		namespace, _ := store.GetValue(fmt.Sprintf("/pods/%s/namespace", pod))
		for _, port := range store.List(fmt.Sprintf("/ports/%s/", pod)) {
			ch <- prometheus.MustNewConstMetric(collector.openPort, prometheus.GaugeValue, 1, name, namespace, port)
		}
	}

	// fetch ScanDuration in seconds  - this is a gauge since is absolute value
	scanDurationSecondsString, err := store.GetValue("/timings/duration")
	if err != nil {
		log.Debugf("No Scan Duration in the store")
		scanDurationSecondsString = "0.0"
	}
	scanDurationSeconds, err := strconv.ParseFloat(scanDurationSecondsString, 64)
	if err != nil {
		log.Debugf("Failed to convert scan duration to float : %f", scanDurationSecondsString)
		scanDurationSeconds = 0
	}

	// fetch last successful scan in epoch - this is a gauge since is an absolute value
	lastSuccessfulScanEpochString, err := store.GetValue("/timings/last")
	if err != nil {
		log.Debugf("No last epoch in the store")
		lastSuccessfulScanEpochString = "0"
	}
	lastSuccessfulScanEpoch, err := strconv.ParseFloat(lastSuccessfulScanEpochString, 64)
	if err != nil {
		log.Debugf("Failed to convert last epoch : %f", lastSuccessfulScanEpochString)
		lastSuccessfulScanEpoch = 0
	}

	ch <- prometheus.MustNewConstMetric(collector.scanDuration, prometheus.GaugeValue, scanDurationSeconds)
	ch <- prometheus.MustNewConstMetric(collector.lastSuccessfulScan, prometheus.GaugeValue, lastSuccessfulScanEpoch)

	// Number of pods

	numberOfScannedPods := store.ListDir("/pods")
	numberOfOpenPods := store.ListDir("/ports")

	ch <- prometheus.MustNewConstMetric(collector.scannedPodsTotal, prometheus.GaugeValue, float64(len(numberOfScannedPods)))
	ch <- prometheus.MustNewConstMetric(collector.openPodsTotal, prometheus.GaugeValue, float64(len(numberOfOpenPods)))

}

func portScanRegister() {
	collector := newPortScanCollector()
	prometheus.MustRegister(collector)
}

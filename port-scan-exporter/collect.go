package main

import (
	"context"
	"fmt"
	"time"

	portscanner "github.com/anvie/port-scanner"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

// Collect List of pods we want to scan and them to Memory Store
func collectPods() {

	log.Infof("Collecting Pods")

	ctx := context.TODO()
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	log.Debugf("There are %d pods in the cluster", len(pods.Items))

	// Clear the DB before evaluating currently running pods
	store.Del("/pods")
	store.Del("/ports")

	// Filter Pods to exclude those running with host network
	for _, p := range pods.Items {
		if p.Spec.HostNetwork {
			log.Debugf("Excluding host network pod %s", p.GetName())
		} else {
			log.Debugf("Adding pod %s to list of pods to scan", p.GetName())
			store.Set(fmt.Sprintf("/pods/%s/ip", p.GetName()), p.Status.PodIP)
			store.Set(fmt.Sprintf("/pods/%s/name", p.GetName()), p.GetName())
			store.Set(fmt.Sprintf("/pods/%s/namespace", p.GetName()), p.GetNamespace())
		}
	}

}

// Scan Pods in Memory Store
func scanPods(min int, max int) {

	pods := store.ListDir("/pods")
	log.Infof("Scanning %d Pods, min port:%d max port:%d", len(pods), min, max)

	// https://stackoverflow.com/questions/25306073/always-have-x-number-of-goroutines-running-at-any-time
	maxScanners := 10
	guard := make(chan int, maxScanners)
	for _, p := range pods {
		guard <- 1 // would block if guard channel is already filled
		go func() {
			name, _ := store.GetValue(fmt.Sprintf("/pods/%s/name", p))
			ip, _ := store.GetValue(fmt.Sprintf("/pods/%s/ip", p))
			log.Debugf("Scanning %s Pod with ip %s", name, ip)

			// scan host with a 2 second timeout per port in 5 concurrent threads
			ps := portscanner.NewPortScanner(ip, 2*time.Second, 5)
			openedPorts := ps.GetOpenedPort(min, max)

			for i := 0; i < len(openedPorts); i++ {
				store.Set(fmt.Sprintf("/ports/%s/%d", name, openedPorts[i]), "true")
			}

			<-guard // removes an int from guard, allowing another to proceed
		}()
	}
}

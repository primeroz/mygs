package main

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	//"github.com/anvie/port-scanner"
)

func collectPods() {

	start := time.Now()

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

	timeTrack(start, "Collecting Pods")

}

func scanPods(min uint64, max uint64) {

	start := time.Now()

	pods := store.ListDir("/pods")
	log.Infof("Scanning %d Pods, min port:%d max port:%d", len(pods), min, max)

	// https://stackoverflow.com/questions/25306073/always-have-x-number-of-goroutines-running-at-any-time
	maxScanners := 10
	guard := make(chan int, maxScanners)
	for _, p := range pods {
		guard <- 1 // would block if guard channel is already filled
		go func() {
			_, name := store.Get(fmt.Sprintf("/pods/%s/name", p))
			_, ip := store.Get(fmt.Sprintf("/pods/%s/ip", p))
			log.Debugf("Scanning %d Pod with ip %s", name, ip)
			<-guard // removes an int from guard, allowing another to proceed
		}()
	}

	timeTrack(start, "Scanning Pods")
}

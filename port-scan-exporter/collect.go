package main

import (
	"context"
	"fmt"
	"time"

	//			v1 "k8s.io/api/apps/v1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func collectPods() {

	start := time.Now()

	ctx := context.TODO()
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	log.Infof("There are %d pods in the cluster", len(pods.Items))

	// Filter Pods to exclude those running with host network
	for _, p := range pods.Items {
		if p.Spec.HostNetwork {
			log.Infof("Excluding host network pod %s", p.GetName())
		} else {
			log.Infof("Adding pod %s to list of pods to scan", p.GetName())
			store.Set(fmt.Sprintf("/pods/%s/ip", p.GetName()), p.Status.PodIP)
			store.Set(fmt.Sprintf("/pods/%s/name", p.GetName()), p.GetName())
			store.Set(fmt.Sprintf("/pods/%s/namespace", p.GetName()), p.GetNamespace())
		}
	}

	// TMP Check
	keys := store.List("/pods/")
	log.Infof("There are %d keys in the store", len(keys))
	for _, key := range keys {
		log.Infof("Key: %s", key)
	}
	store.Del("/pods")
	testkeys := store.List("/")
	log.Infof("There are %d keys in the store", len(testkeys))

	timeTrack(start, "Collecting Pods")

}

func scanPods() {

	start := time.Now()

	pods, err := store.GetAll("/*")
	if err != nil {
		panic(err.Error())
	}
	log.Infof("Scanning %s Pods", len(pods))

	timeTrack(start, "Collecting Pods")
}

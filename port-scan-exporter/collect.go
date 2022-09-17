package main

import (
	"context"
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
		}
		//else {
		//  store.Set(fmt.Printf("/pods/%s/ip", p.GetName()), p.Status.PodIP)
		//}
	}

	timeTrack(start, "Collecting Pods")

}

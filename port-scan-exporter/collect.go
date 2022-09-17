package main

import (
	"context"

	//		"fmt"
	//
	//			v1 "k8s.io/api/apps/v1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func collectPods() {

	ctx := context.TODO()
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	log.Infof("There are %d pods in the cluster\n", len(pods.Items))
}

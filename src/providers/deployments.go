package providers

import (
	"context"
	"github.com/imkira/go-observer"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"log"
)

type DeploymentProvider struct {
	InitialDeployments []v1.Deployment
	ResourceVersion    string
	Observer           *WatchObserver
}

// New Deployments Provider
func NewDeploymentProvider(k8s *kubernetes.Clientset) *DeploymentProvider {
	ctx := context.Background()
	timeout := int64(60)
	initial, err := k8s.AppsV1().Deployments("airflow").List(ctx, metav1.ListOptions{
		TimeoutSeconds: &timeout,
	})

	if err != nil {
		log.Fatal(err)
	}

	dp := &DeploymentProvider{
		InitialDeployments: initial.Items,
		ResourceVersion:    initial.ResourceVersion,
	}

	watcher, err := k8s.AppsV1().
		Deployments("airflow").
		Watch(ctx, metav1.ListOptions{ResourceVersion: dp.ResourceVersion})

	if err != nil {
		log.Fatal(err)
	}

	// Create a property to base observers on
	prop := observer.NewProperty(watch.Event{})

	wo := &WatchObserver{
		EventChannel:  watcher.ResultChan(),
		EventProperty: prop,
	}
	wo.Run()

	dp.Observer = wo
	return dp
}

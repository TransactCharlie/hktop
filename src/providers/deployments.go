package providers

import (
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/watch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/imkira/go-observer"
	"k8s.io/client-go/kubernetes"
	"log"
)

type DeploymentProvider struct {
	InitialDeployments []v1beta1.Deployment
	ResourceVersion string
	DeploymentObserver *WatchObserver
}

// New Deployments Provider
func NewDeploymentProvider(k8s *kubernetes.Clientset) *DeploymentProvider {
	initial, err := k8s.ExtensionsV1beta1().Deployments("").List(metav1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	}

	dp := &DeploymentProvider{
		InitialDeployments: initial.Items,
		ResourceVersion: initial.ResourceVersion,
	}

	watcher, err := k8s.ExtensionsV1beta1().
		Deployments("").
		Watch(metav1.ListOptions{ResourceVersion: dp.ResourceVersion})

	if err != nil {
		log.Fatal(err)
	}

	// Create a property to base observers on
	prop := observer.NewProperty(watch.Event{})

	wo := &WatchObserver{
		EventChannel: watcher.ResultChan(),
		EventProperty: prop,
	}
	wo.Run()

	dp.DeploymentObserver = wo
	return dp
}

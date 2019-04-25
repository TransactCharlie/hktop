package providers

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/imkira/go-observer"
	"k8s.io/client-go/kubernetes"
	"log"
)

type ServiceProvider struct {
	InitialServices []v1.Service
	ResourceVersion string
	ServiceObserver *WatchObserver
}

// New Service Provider
func NewServiceProvider(k8s *kubernetes.Clientset) *ServiceProvider {

	initial, err := k8s.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	sp := &ServiceProvider{
		InitialServices: initial.Items,
		ResourceVersion: initial.ResourceVersion,
	}

	// Watch Services
	watcher, err := k8s.CoreV1().
		Services("").
		Watch(metav1.ListOptions{ResourceVersion:sp.ResourceVersion})

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

	sp.ServiceObserver = wo
	return sp
}
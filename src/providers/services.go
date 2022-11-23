package providers

import (
	"context"
	"github.com/imkira/go-observer"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"log"
)

type ServiceProvider struct {
	InitialServices []v1.Service
	ResourceVersion string
	Observer        *WatchObserver
}

// New Service Provider
func NewServiceProvider(k8s *kubernetes.Clientset) *ServiceProvider {
	ctx := context.Background()
	initial, err := k8s.CoreV1().Services("").List(ctx, metav1.ListOptions{})
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
		Watch(ctx, metav1.ListOptions{ResourceVersion: sp.ResourceVersion})

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

	sp.Observer = wo
	return sp
}

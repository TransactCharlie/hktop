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

type PodProvider struct {
	InitialPods     []v1.Pod
	ResourceVersion string
	Observer        *WatchObserver
}

// Creates a (running) PodProvider with an initialised k8s pod watch and list
func NewPodProvider(k8s *kubernetes.Clientset) *PodProvider {
	ctx := context.Background()
	initialPods, err := k8s.CoreV1().Pods("airflow").List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	pp := &PodProvider{
		InitialPods:     initialPods.Items,
		ResourceVersion: initialPods.ResourceVersion,
	}

	watcher, err := k8s.CoreV1().Pods("").Watch(ctx, metav1.ListOptions{ResourceVersion: pp.ResourceVersion})
	if err != nil {
		log.Fatal(err)
	}

	// Create a property to base observers on
	prop := observer.NewProperty(watch.Event{})

	podObserver := &WatchObserver{
		EventChannel:  watcher.ResultChan(),
		EventProperty: prop,
	}
	podObserver.Run()

	pp.Observer = podObserver
	return pp
}

package providers

import (
	"k8s.io/apimachinery/pkg/watch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/imkira/go-observer"
	"k8s.io/client-go/kubernetes"
	"log"
)

// New Deployments Observer
func NewServiceObserver(k8s *kubernetes.Clientset) *WatchObserver {
	watcher, err := k8s.CoreV1().Services("").Watch(metav1.ListOptions{})
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
	return wo
}

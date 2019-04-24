package providers

import (
	"k8s.io/apimachinery/pkg/watch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/imkira/go-observer"
	"k8s.io/client-go/kubernetes"
	"log"
)

// Creates a (running) new NodeObserver with an initialised k8s node watch
func NewNodeObserver(clientset *kubernetes.Clientset) *WatchObserver {
	watcher, err := clientset.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Create a property to base observers on
	prop := observer.NewProperty(watch.Event{})

	no := &WatchObserver{
		EventChannel: watcher.ResultChan(),
		EventProperty: prop,
	}
	no.Run()
	return no
}

package providers

import (
	"k8s.io/apimachinery/pkg/watch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/imkira/go-observer"
	"k8s.io/client-go/kubernetes"
	"log"
)

// New Node Observer
func NewNodeObserver(k8s *kubernetes.Clientset) *WatchObserver {
	watcher, err := k8s.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	no := &WatchObserver{
		EventChannel: watcher.ResultChan(),
		EventProperty: observer.NewProperty(watch.Event{}),
	}
	no.Run()
	return no
}

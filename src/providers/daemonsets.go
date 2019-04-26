package providers

import (
	"github.com/imkira/go-observer"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"log"
)

type DaemonSetProvider struct {
	InitialDaemonSets []v1beta1.DaemonSet
	ResourceVersion   string
	Observer          *WatchObserver
}

// New DaemonSet Provider
func NewDaemonSetProvider(k8s *kubernetes.Clientset) *DaemonSetProvider {
	initial, err := k8s.ExtensionsV1beta1().
		DaemonSets("").
		List(metav1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	}

	dp := &DaemonSetProvider{
		InitialDaemonSets: initial.Items,
		ResourceVersion:   initial.ResourceVersion,
	}

	watcher, err := k8s.ExtensionsV1beta1().
		DaemonSets("").
		Watch(metav1.ListOptions{ResourceVersion: dp.ResourceVersion})

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

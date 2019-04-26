package providers

import (
	"github.com/imkira/go-observer"
	"k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"log"
)

type StatefulSetProvider struct {
	InitialStatefulSets []v1beta1.StatefulSet
	ResourceVersion     string
	Observer            *WatchObserver
}

// New StatefulSetProvider Provider
func NewStatefulSetProvider(k8s *kubernetes.Clientset) *StatefulSetProvider {
	initial, err := k8s.AppsV1beta1().
		StatefulSets("").
		List(metav1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	}

	provider := &StatefulSetProvider{
		InitialStatefulSets: initial.Items,
		ResourceVersion:     initial.ResourceVersion,
	}

	watcher, err := k8s.AppsV1beta1().
		StatefulSets("").
		Watch(metav1.ListOptions{ResourceVersion: provider.ResourceVersion})

	if err != nil {
		log.Fatal(err)
	}

	observer := &WatchObserver{
		EventChannel:  watcher.ResultChan(),
		EventProperty: observer.NewProperty(watch.Event{}),
	}
	observer.Run()

	provider.Observer = observer
	return provider
}

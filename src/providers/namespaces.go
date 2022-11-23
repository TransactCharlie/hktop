package providers

import (
	"context"
	"github.com/imkira/go-observer"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"log"
)

type NamespaceProvider struct {
	InitialNamespaces []v1.Namespace
	ResourceVersion   string
	Observer          *WatchObserver
}

// New PersistentVolume Provider
func NewNamespaceProvider(k8s *kubernetes.Clientset) *NamespaceProvider {
	ctx := context.Background()
	initial, err := k8s.CoreV1().
		Namespaces().
		List(ctx, metav1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	}

	provider := &NamespaceProvider{
		InitialNamespaces: initial.Items,
		ResourceVersion:   initial.ResourceVersion,
	}

	watcher, err := k8s.CoreV1().
		Namespaces().
		Watch(ctx, metav1.ListOptions{ResourceVersion: provider.ResourceVersion})

	if err != nil {
		log.Fatal(err)
	}

	wo := &WatchObserver{
		EventChannel:  watcher.ResultChan(),
		EventProperty: observer.NewProperty(watch.Event{}),
	}
	wo.Run()

	provider.Observer = wo
	return provider
}

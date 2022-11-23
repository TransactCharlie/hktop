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

type NodeProvider struct {
	InitialNodes    []v1.Node
	ResourceVersion string
	Observer        *WatchObserver
}

// New Node Provider
func NewNodeProvider(k8s *kubernetes.Clientset) *NodeProvider {
	ctx := context.Background()

	initialNodes, err := k8s.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	watcher, err := k8s.CoreV1().Nodes().Watch(ctx, metav1.ListOptions{ResourceVersion: initialNodes.ListMeta.ResourceVersion})
	if err != nil {
		log.Fatal(err)
	}

	no := &WatchObserver{
		EventChannel:  watcher.ResultChan(),
		EventProperty: observer.NewProperty(watch.Event{}),
	}
	no.Run()

	np := &NodeProvider{
		InitialNodes:    initialNodes.Items,
		ResourceVersion: initialNodes.ListMeta.ResourceVersion,
		Observer:        no,
	}

	return np
}

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

type PersistentVolumeProvider struct {
	InitialPersistentVolumes []v1.PersistentVolume
	ResourceVersion          string
	Observer                 *WatchObserver
}

// New PersistentVolumeProvider
func NewPersistentVolumeProvider(k8s *kubernetes.Clientset) *PersistentVolumeProvider {
	ctx := context.Background()
	initial, err := k8s.CoreV1().
		PersistentVolumes().
		List(ctx, metav1.ListOptions{})

	if err != nil {
		log.Fatal(err)
	}

	pp := &PersistentVolumeProvider{
		InitialPersistentVolumes: initial.Items,
		ResourceVersion:          initial.ResourceVersion,
	}

	watcher, err := k8s.CoreV1().
		PersistentVolumes().
		Watch(ctx, metav1.ListOptions{ResourceVersion: pp.ResourceVersion})

	if err != nil {
		log.Fatal(err)
	}

	wo := &WatchObserver{
		EventChannel:  watcher.ResultChan(),
		EventProperty: observer.NewProperty(watch.Event{}),
	}
	wo.Run()

	pp.Observer = wo
	return pp
}

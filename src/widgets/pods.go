package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	"github.com/imkira/go-observer"
	p "github.com/transactcharlie/hktop/src/providers"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)


type KubernetesPods struct {
	*ui.List
	Events observer.Stream
	stop chan bool
}

func NewKubernetesPods(observer *p.WatchObserver) *KubernetesPods {
	kn := &KubernetesPods{
		List:       ui.NewList(),
		Events: observer.RegisterObserver(),
		stop:       make(chan bool),
	}
	kn.Rows = []string{}
	kn.Title = "K8S Pods"
	go func() { _ = kn.Update() }()
	kn.Run()
	return kn
}

func (kn *KubernetesPods) Run() {
	go func() {
		for {
			select {
			case <-kn.stop:
				return
			case <- kn.Events.Changes():
				kn.Events.Next()
				event := kn.Events.Value().(watch.Event)
				pod, _ := event.Object.(*v1.Pod)
				switch event.Type {
				case watch.Added:
					kn.Rows = append(kn.Rows, pod.Name)
				case watch.Deleted:
					for ix, val := range kn.Rows {
						if val == pod.Name {
							kn.Rows = append(kn.Rows[:ix], kn.Rows[ix+1:]...)
							break
						}
					}
				case watch.Modified:
					continue
				}
			}
		}
	}()
}

func (kn *KubernetesPods) Stop() bool {
	kn.stop <- true
	return true
}

func (kn *KubernetesPods) Update() error {
	return nil
}

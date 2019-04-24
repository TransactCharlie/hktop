package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	"github.com/imkira/go-observer"
	p "github.com/transactcharlie/hktop/src/providers"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type KubernetesNodes struct {
	*ui.List
	Events observer.Stream
	stop chan bool
}

func NewKubernetesNode(observer *p.WatchObserver) *KubernetesNodes {
	kn := &KubernetesNodes{
		List: ui.NewList(),
		Events: observer.RegisterObserver(),
		stop: make(chan bool),
	}
	kn.Rows = []string{}
	kn.Title = "K8S Nodes"
	go func() {_ = kn.Update()}()
	kn.Run()
	return kn
}

func (kn *KubernetesNodes) Run() {
	go func() {
		for {
			select {
			case <- kn.stop:
				return
			// Deal with a change
			case <- kn.Events.Changes():
				// advance to new value
				kn.Events.Next()
				event := kn.Events.Value().(watch.Event)
				node, _ := event.Object.(*v1.Node)
				switch event.Type {
				case watch.Added:
					kn.Rows = append(kn.Rows, node.Name)
				}
			}
		}
	}()
}

func (kn *KubernetesNodes) Stop() bool {
	kn.stop <- true
	return true
}

func (kn *KubernetesNodes) Update() error {
	return nil
}
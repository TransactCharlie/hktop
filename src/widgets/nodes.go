package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	"github.com/imkira/go-observer"
	p "github.com/transactcharlie/hktop/src/providers"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type NodeListWidget struct {
	*ui.List
	Events observer.Stream
	stop chan bool
}

func NewNodeListWidget(observer *p.WatchObserver) *NodeListWidget {
	nlw := &NodeListWidget{
		List: ui.NewList(),
		Events: observer.RegisterObserver(),
		stop: make(chan bool),
	}
	nlw.Rows = []string{}
	nlw.Title = "K8S Nodes"
	nlw.Run()
	return nlw
}

func (nlw *NodeListWidget) Run() {
	go func() {
		for {
			select {
			case <- nlw.stop:
				return
			// Deal with a change
			case <- nlw.Events.Changes():
				// advance to new value
				nlw.Events.Next()
				event := nlw.Events.Value().(watch.Event)
				node, _ := event.Object.(*v1.Node)
				switch event.Type {
				case watch.Added:
					nlw.Rows = append(nlw.Rows, node.Name)
				}
			}
		}
	}()
}

func (nlw *NodeListWidget) Stop() bool {
	nlw.stop <- true
	return true
}

func (nlw *NodeListWidget) Update() error {
	return nil
}
package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	observer "github.com/imkira/go-observer"
	"github.com/transactcharlie/hktop/src/providers"
	"k8s.io/apimachinery/pkg/watch"
	"fmt"
)

type SummaryWidget struct {
	*ui.Table
	NodeEvents observer.Stream
	PodEvents observer.Stream
	DeploymentEvents observer.Stream
	ServiceEvents observer.Stream
	stop chan bool
	nodeCount int
	podCount int
	deploymentCount int
	serviceCount int
}

func NewSummaryWidget(no, po, do, so *providers.WatchObserver) *SummaryWidget {
	sw := &SummaryWidget{
		Table: ui.NewTable(),
		NodeEvents: no.RegisterObserver(),
		PodEvents: po.RegisterObserver(),
		DeploymentEvents: do.RegisterObserver(),
		ServiceEvents: so.RegisterObserver(),
		stop: make(chan bool),
	}
	sw.Rows = [][]string{
		{"Nodes", ""},
		{"Pods", ""},
		{"Deployments", ""},
		{"Services", ""},
	}
	sw.Title = "Summary"
	sw.RowSeparator = false
	sw.Run()
	return sw
}

func (sw *SummaryWidget) Stop() bool {
	sw.stop <- true
	return true
}

func (sw *SummaryWidget) Update() error {
	sw.Rows[0][1] = fmt.Sprintf("%v", sw.nodeCount)
	sw.Rows[1][1] = fmt.Sprintf("%v", sw.podCount)
	sw.Rows[2][1] = fmt.Sprintf("%v", sw.deploymentCount)
	sw.Rows[3][1] = fmt.Sprintf("%v", sw.serviceCount)
	return nil
}

func (sw *SummaryWidget) Run() {
	go func() {
		for {
			select {
			case <- sw.stop:
				return
			case <- sw.NodeEvents.Changes():
				sw.processNodeChange()
			case <- sw.PodEvents.Changes():
				sw.processPodChange()
			case <- sw.DeploymentEvents.Changes():
				sw.processDeploymentChange()
			case <- sw.ServiceEvents.Changes():
				sw.processServiceChange()
			}
		}
	}()
}

func (sw *SummaryWidget) processNodeChange() {
	sw.NodeEvents.Next()
	event := sw.NodeEvents.Value().(watch.Event)
	switch event.Type {
	case watch.Added:
		sw.nodeCount += 1
	case watch.Deleted:
		sw.nodeCount -= 1
	}
	_ = sw.Update()
}

func (sw *SummaryWidget) processPodChange() {
	sw.PodEvents.Next()
	event := sw.PodEvents.Value().(watch.Event)
	switch event.Type {
	case watch.Added:
		sw.podCount += 1
	case watch.Deleted:
		sw.podCount -= 1
	}
	_ = sw.Update()
}

func (sw *SummaryWidget) processDeploymentChange() {
	sw.DeploymentEvents.Next()
	event := sw.DeploymentEvents.Value().(watch.Event)
	switch event.Type {
	case watch.Added:
		sw.deploymentCount += 1
	case watch.Deleted:
		sw.deploymentCount -= 1
	}
	_ = sw.Update()
}

func (sw *SummaryWidget) processServiceChange() {
	sw.ServiceEvents.Next()
	event := sw.ServiceEvents.Value().(watch.Event)
	switch event.Type {
	case watch.Added:
		sw.serviceCount += 1
	case watch.Deleted:
		sw.serviceCount -= 1
	}
	_ = sw.Update()
}
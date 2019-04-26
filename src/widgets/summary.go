package widgets

import (
	"fmt"
	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3/widgets"
	observer "github.com/imkira/go-observer"
	p "github.com/transactcharlie/hktop/src/providers"
	"k8s.io/apimachinery/pkg/watch"
)

type SummaryWidget struct {
	*ui.Table
	NodeEvents             observer.Stream
	PodEvents              observer.Stream
	DeploymentEvents       observer.Stream
	ServiceEvents          observer.Stream
	DaemonSetEvents        observer.Stream
	PersistentVolumeEvents observer.Stream
	NamespaceEvents        observer.Stream
	StatefulSetEvents      observer.Stream
	stop                   chan bool
	nodeCount              int
	podCount               int
	deploymentCount        int
	serviceCount           int
	daemonSetCount         int
	persistentVolumeCount  int
	namespaceCount         int
	statefulSetCount       int
}

func NewSummaryWidget(np *p.NodeProvider,
	pp *p.PodProvider,
	dp *p.DeploymentProvider,
	sp *p.ServiceProvider,
	dsp *p.DaemonSetProvider,
	pvp *p.PersistentVolumeProvider,
	nsp *p.NamespaceProvider,
	ssp *p.StatefulSetProvider,
) *SummaryWidget {
	sw := &SummaryWidget{
		Table:                  ui.NewTable(),
		NodeEvents:             np.NodeObserver.RegisterObserver(),
		PodEvents:              pp.PodObserver.RegisterObserver(),
		DeploymentEvents:       dp.DeploymentObserver.RegisterObserver(),
		ServiceEvents:          sp.ServiceObserver.RegisterObserver(),
		DaemonSetEvents:        dsp.Observer.RegisterObserver(),
		PersistentVolumeEvents: pvp.Observer.RegisterObserver(),
		NamespaceEvents:        nsp.Observer.RegisterObserver(),
		StatefulSetEvents:      ssp.Observer.RegisterObserver(),
		stop:                   make(chan bool),
		nodeCount:              len(np.InitialNodes),
		podCount:               len(pp.InitialPods),
		deploymentCount:        len(dp.InitialDeployments),
		serviceCount:           len(sp.InitialServices),
		daemonSetCount:         len(dsp.InitialDaemonSets),
		persistentVolumeCount:  len(pvp.InitialPersistentVolumes),
		namespaceCount:         len(nsp.InitialNamespaces),
		statefulSetCount:       len(ssp.InitialStatefulSets),
	}
	sw.Rows = [][]string{
		{"Nodes", ""},
		{"Pods", ""},
		{"Deployments", ""},
		{"Services", ""},
		{"Daemon Sets", ""},
		{"Persistent Volumes", ""},
		{"Namespaces", ""},
		{"Stateful Sets", ""},
	}
	sw.FillRow = false
	sw.TextAlignment = termui.AlignLeft
	sw.ColumnResizer = func() {
		width := sw.Inner.Dx()
		countWidth := 0
		for _, r := range sw.Rows {
			if len(r[1]) > countWidth {
				countWidth = len(r[1])
			}
		}
		// pad countwidth
		countWidth += 2
		textWidth := width - countWidth
		sw.ColumnWidths = []int{textWidth, countWidth}
		return
	}
	sw.Title = "Summary"
	sw.RowSeparator = false
	sw.Run()
	go func() { _ = sw.Update() }()
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
	sw.Rows[4][1] = fmt.Sprintf("%v", sw.daemonSetCount)
	sw.Rows[5][1] = fmt.Sprintf("%v", sw.persistentVolumeCount)
	sw.Rows[6][1] = fmt.Sprintf("%v", sw.namespaceCount)
	sw.Rows[7][1] = fmt.Sprintf("%v", sw.statefulSetCount)
	return nil
}

func (sw *SummaryWidget) Run() {
	go func() {
		for {
			select {
			case <-sw.stop:
				return
			case <-sw.NodeEvents.Changes():
				sw.processNodeChange()
			case <-sw.PodEvents.Changes():
				sw.processPodChange()
			case <-sw.DeploymentEvents.Changes():
				sw.processDeploymentChange()
			case <-sw.ServiceEvents.Changes():
				sw.processServiceChange()
			case <-sw.DaemonSetEvents.Changes():
				sw.processDaemonSetChange()
			case <-sw.PersistentVolumeEvents.Changes():
				sw.processPersistentVolumeChange()
			case <-sw.NamespaceEvents.Changes():
				sw.processNamespaceChange()
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

func (sw *SummaryWidget) processDaemonSetChange() {
	sw.DaemonSetEvents.Next()
	event := sw.DaemonSetEvents.Value().(watch.Event)
	switch event.Type {
	case watch.Added:
		sw.daemonSetCount += 1
	case watch.Deleted:
		sw.daemonSetCount -= 1
	}
	_ = sw.Update()
}

func (sw *SummaryWidget) processPersistentVolumeChange() {
	sw.PersistentVolumeEvents.Next()
	event := sw.PersistentVolumeEvents.Value().(watch.Event)
	switch event.Type {
	case watch.Added:
		sw.persistentVolumeCount += 1
	case watch.Deleted:
		sw.persistentVolumeCount -= 1
	}
	_ = sw.Update()
}

func (sw *SummaryWidget) processNamespaceChange() {
	sw.NamespaceEvents.Next()
	event := sw.NamespaceEvents.Value().(watch.Event)
	switch event.Type {
	case watch.Added:
		sw.namespaceCount += 1
	case watch.Deleted:
		sw.namespaceCount -= 1
	}
	_ = sw.Update()
}

func (sw *SummaryWidget) statefulSetChange() {
	sw.StatefulSetEvents.Next()
	event := sw.StatefulSetEvents.Value().(watch.Event)
	switch event.Type {
	case watch.Added:
		sw.statefulSetCount += 1
	case watch.Deleted:
		sw.statefulSetCount -= 1
	}
	_ = sw.Update()
}

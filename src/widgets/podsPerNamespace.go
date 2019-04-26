package widgets

import (
	"fmt"
	ui "github.com/gizak/termui/v3/widgets"
	"github.com/imkira/go-observer"
	p "github.com/transactcharlie/hktop/src/providers"
	"k8s.io/api/core/v1"
	"sort"
)

type PodsPerNamespaceWidget struct {
	*ui.PieChart
	PodEvents        observer.Stream
	NamespaceEvents  observer.Stream
	stop             chan bool
	namespacePods    map[string]map[string]v1.Pod
	sortedNamespaces []string
}

func NewPodsPerNamespaceWidget(pp *p.PodProvider, nsp *p.NamespaceProvider) *PodsPerNamespaceWidget {
	w := &PodsPerNamespaceWidget{
		PieChart:         ui.NewPieChart(),
		PodEvents:        pp.PodObserver.RegisterObserver(),
		NamespaceEvents:  nsp.Observer.RegisterObserver(),
		stop:             make(chan bool),
		namespacePods:    make(map[string]map[string]v1.Pod),
		sortedNamespaces: []string{},
	}
	w.Title = "Pods per Namespace"
	w.PaddingTop = 1
	w.PaddingBottom = 2
	w.AngleOffset = 5

	for _, namespace := range nsp.InitialNamespaces {
		w.namespacePods[namespace.Name] = make(map[string]v1.Pod)
		w.sortedNamespaces = append(w.sortedNamespaces, namespace.Name)
	}
	sort.Strings(w.sortedNamespaces)

	for _, pod := range pp.InitialPods {
		w.namespacePods[pod.Namespace][pod.Name] = pod
	}

	// Initialise Data
	w.Data = []float64{}
	for _, namespace := range w.sortedNamespaces {
		w.Data = append(w.Data, float64(len(w.namespacePods[namespace])))
	}

	// Hook to return labels for the pie chart
	w.LabelFormatter = w.pieLabels
	return w
}

func (w *PodsPerNamespaceWidget) pieLabels(i int, v float64) string {
	if v >= 1.0 {
		return fmt.Sprintf("%s: %v", w.sortedNamespaces[i], v)
	}
	return ""
}

func (w *PodsPerNamespaceWidget) Run() {
	return
}

func (w *PodsPerNamespaceWidget) Update() error {

	return nil
}

func (w *PodsPerNamespaceWidget) Stop() bool {
	return true
}

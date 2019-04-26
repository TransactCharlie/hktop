package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	"github.com/imkira/go-observer"
	p "github.com/transactcharlie/hktop/src/providers"
)

type PodsPerNamespaceWidget struct {
	*ui.PieChart
	Events observer.Stream
	stop   chan bool
}

func NewPodsPerNamespaceWidget(pp *p.PodProvider) *PodsPerNamespaceWidget {
	w := &PodsPerNamespaceWidget{
		PieChart: ui.NewPieChart(),
		Events:   pp.PodObserver.RegisterObserver(),
		stop:     make(chan bool),
	}
	w.Data = []float64{1, 1, 1, 1}
	return w
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

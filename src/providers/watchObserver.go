package providers

import (
	"k8s.io/apimachinery/pkg/watch"
	"github.com/imkira/go-observer"
)

type WatchObserver struct {
	stop chan bool
	EventChannel <-chan watch.Event
	EventProperty observer.Property
}

// Register a new Observer
func (wo *WatchObserver) RegisterObserver() observer.Stream {
	return wo.EventProperty.Observe()
}

// Runs the nodeObserver processing the watch for nodes and
// publishing events to observers
func (wo *WatchObserver) Run() {
	go func() {
		for {
			select {
			case <-wo.stop:
				return
			case event := <-wo.EventChannel:
				wo.EventProperty.Update(event)
			}
		}
	}()
}

// Stops the Observer
func (wo *WatchObserver) Stop() bool {
	wo.stop <- true
	return true
}
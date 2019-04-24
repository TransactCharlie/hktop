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
func (no *WatchObserver) RegisterObserver() observer.Stream {
	return no.EventProperty.Observe()
}

// Runs the nodeObserver processing the watch for nodes and
// publishing events to observers
func (no *WatchObserver) Run() {
	go func() {
		for {
			select {
			case <-no.stop:
				return
			case event := <-no.EventChannel:
				no.EventProperty.Update(event)
			}
		}
	}()
}

// Stops the Observer
func (no *WatchObserver) Stop() bool {
	no.stop <- true
	return true
}
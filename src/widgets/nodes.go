package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	watch "k8s.io/apimachinery/pkg/watch"
	"log"
	"time"
)

type KubernetesNodes struct {
	*ui.List
	clientset *kubernetes.Clientset
	updateTick <-chan time.Time
	nodeWatch <-chan watch.Event
	stop chan bool
}

func NewKubernetesNode(clientset *kubernetes.Clientset) *KubernetesNodes {
	kn := &KubernetesNodes{
		List: ui.NewList(),
		clientset: clientset,
		updateTick: time.NewTicker(time.Second * 10).C,
		nodeWatch: createNodeWatch(clientset),
		stop: make(chan bool),
	}
	kn.Rows = []string{}
	kn.Title = "K8S Nodes"
	go func() {_ = kn.Update()}()
	kn.Run()
	return kn
}

func createNodeWatch(clientset *kubernetes.Clientset) <-chan watch.Event {
	watcher, err := clientset.CoreV1().Nodes().Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return watcher.ResultChan()
}

func (kn *KubernetesNodes) Run() {
	go func() {
		for {
			select {
			case <- kn.stop:
				return
			case <- kn.updateTick:
				_ = kn.Update()
			case _ = <- kn.nodeWatch:
				continue
			}
		}
	}()
}

func (kn *KubernetesNodes) Stop() bool {
	kn.stop <- true
	return true
}

func (kn *KubernetesNodes) Update() error {
	nodes, err := kn.K8SNodes()
	if err != nil {
		return err
	}
	nodeDetails := nodes.Items
	newRows := []string{}
	for _, nd := range(nodeDetails) {
		newRows = append(newRows, nd.Name)
	}
	kn.Rows = newRows
	return nil
}

func (kn *KubernetesNodes) K8SNodes() (*v1.NodeList, error) {
	nodes, err := kn.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	return nodes, err
}
package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"time"
)

type KubernetesPods struct {
	*ui.List
	clientset *kubernetes.Clientset
	updateTick <-chan time.Time
	podWatch <- chan watch.Event
	stop chan bool
}

func (kn *KubernetesPods) Run() {
	go func() {
		for {
			select {
			case <- kn.stop:
				return
			case event := <- kn.podWatch:
				pod, ok := event.Object.(*v1.Pod)
				if !ok {
					log.Fatal("unexpected type")
				}
				switch event.Type {
				case watch.Added:
					kn.Rows = append(kn.Rows, pod.Name)
				}
			}
		}
	}()

}

func (kn *KubernetesPods) Stop() bool {
	kn.stop <- true
	return true
}

func createPodWatch(clientset *kubernetes.Clientset) <-chan watch.Event {
	watcher, err := clientset.CoreV1().Pods("").Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	return watcher.ResultChan()
}

func NewKubernetesPods(clientset *kubernetes.Clientset) *KubernetesPods {
	kn := &KubernetesPods{
		List: ui.NewList(),
		clientset: clientset,
		updateTick: time.NewTicker(time.Second * 10).C,
		podWatch: createPodWatch(clientset),
		stop: make(chan bool),
	}
	kn.Rows = []string{}
	kn.Title = "K8S Pods"
	go func() {_ = kn.Update()}()
	kn.Run()
	return kn
}

func (kn *KubernetesPods) K8SPods() (*v1.PodList, error) {
	pods, err := kn.clientset.CoreV1().Pods("kube-system").List(metav1.ListOptions{})
	return pods, err
}

func (kn *KubernetesPods) Update() error {
	pods, err := kn.K8SPods()
	if err != nil {
		return err
	}
	podDetails := pods.Items
	newPods := []string{}
	for _, pod := range(podDetails) {
		newPods = append(newPods, pod.Name)
	}
	kn.Rows = newPods
	return nil
}
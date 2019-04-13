package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubernetesPods struct {
	*ui.List
	clientset *kubernetes.Clientset
}

func NewKubernetesPods(clientset *kubernetes.Clientset) *KubernetesPods {
	kn := &KubernetesPods{List: ui.NewList(), clientset: clientset}
	kn.Rows = []string{}
	kn.Title = "K8S Pods"

	return kn
}

func (kn *KubernetesPods) K8SPods() (*v1.PodList, error) {
	pods, err := kn.clientset.CoreV1().Pods("kube-system").List(metav1.ListOptions{})
	return pods, err
}

func (kn *KubernetesPods) UpdatePodsList() error {
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
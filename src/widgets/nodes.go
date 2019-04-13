package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubernetesNodes struct {
	*ui.List
	clientset *kubernetes.Clientset
}

func NewKubernetesNode(clientset *kubernetes.Clientset) *KubernetesNodes {
	kn := &KubernetesNodes{List: ui.NewList(), clientset: clientset}
	kn.Rows = []string{}
	kn.Title = "K8S Nodes"

	return kn
}

func (kn *KubernetesNodes) K8SNodes() (*v1.NodeList, error) {
	nodes, err := kn.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	return nodes, err
}

func (kn *KubernetesNodes) UpdateNodeList() error {
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
package main

import (
	"flag"
	ui "github.com/gizak/termui/v3"
	w "github.com/transactcharlie/hktop/src/widgets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"log"
)

var (
	grid *ui.Grid
	exampleParagraphWidget *w.ExampleParagraph
	examplePieWidget *w.ExamplePie
	k8sNodesWidget *w.KubernetesNodes
	k8sPodsWidget *w.KubernetesPods
	updateInterval = time.Second
	clientset *kubernetes.Clientset
)

func initWidgets() {
	exampleParagraphWidget = w.NewExampleParagraph()
	examplePieWidget = w.NewExamplePie()
	k8sNodesWidget = w.NewKubernetesNode(clientset)
	k8sPodsWidget = w.NewKubernetesPods(clientset)
}


func setupGrid() {
	grid = ui.NewGrid()
	grid.Set(
		ui.NewRow(1.0/2, exampleParagraphWidget),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, k8sNodesWidget),
			ui.NewCol(1.0/2, k8sPodsWidget),
		),
	)
}


func eventLoop() {
	drawTicker := time.NewTicker(updateInterval).C
	stateTicker := time.NewTicker(updateInterval).C

	// handles kill signal
	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := ui.PollEvents()

	for {
		select {
		case <-sigTerm:
			return
		case <-stateTicker:
			_ = k8sNodesWidget.UpdateNodeList()
			_ = k8sPodsWidget.UpdatePodsList()
		case <-drawTicker:
			ui.Render(grid)
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render()
			}
		}
	}
}

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}


	defer ui.Close()
	initWidgets()
	setupGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	ui.Render(grid)
	eventLoop()
}
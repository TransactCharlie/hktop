package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	ui "github.com/gizak/termui/v3"
	p "github.com/transactcharlie/hktop/src/providers"
	w "github.com/transactcharlie/hktop/src/widgets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	grid                   *ui.Grid
	exampleParagraphWidget w.K8SWidget
	summaryWidget          w.K8SWidget
	nodeListWidget         w.K8SWidget
	podListWidget          w.K8SWidget
	updateInterval         time.Duration = time.Second
	k8sClientSet           *kubernetes.Clientset
	nodeObserver           *p.WatchObserver
	podObserver            *p.WatchObserver
	deploymentObserver     *p.WatchObserver
	serviceObserver        *p.WatchObserver
)

func initWidgets() {
	exampleParagraphWidget = w.NewExampleParagraph()
	nodeListWidget = w.NewNodeListWidget(nodeObserver)
	podListWidget = w.NewKubernetesPods(podObserver)
	summaryWidget = w.NewSummaryWidget(
										nodeObserver,
										podObserver,
										deploymentObserver,
										serviceObserver,
										)
}

func initObservers() {
	nodeObserver = p.NewNodeObserver(k8sClientSet)
	podObserver = p.NewPodObserver(k8sClientSet)
	deploymentObserver = p.NewDeploymentObserver(k8sClientSet)
	serviceObserver = p.NewServiceObserver(k8sClientSet)
}

func setupGrid() {
	grid = ui.NewGrid()
	grid.Set(
		ui.NewRow(0.25,
			ui.NewCol(0.2, summaryWidget),
			ui.NewCol(0.8, exampleParagraphWidget),
		),
		ui.NewRow(0.75,
			ui.NewCol(0.4, nodeListWidget),
			ui.NewCol(0.6, podListWidget),
		),
	)
}

func eventLoop() {
	drawTicker := time.NewTicker(updateInterval).C

	// handles kill signal
	sigTerm := make(chan os.Signal, 2)
	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)

	uiEvents := ui.PollEvents()

	for {
		select {
		case <-sigTerm:
			return
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
				ui.Render(grid)
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
	k8sClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()
	initObservers()
	initWidgets()
	setupGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
	ui.Render(grid)
	eventLoop()
}

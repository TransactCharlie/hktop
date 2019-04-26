package main

import (
	"flag"
	"k8s.io/client-go/util/homedir"
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
)

var (
	grid                     *ui.Grid
	exampleParagraphWidget   w.K8SWidget
	summaryWidget            w.K8SWidget
	nodeListWidget           w.K8SWidget
	podListWidget            w.K8SWidget
	updateInterval           = time.Second
	nodeProvider             *p.NodeProvider
	podProvider              *p.PodProvider
	deploymentProvider       *p.DeploymentProvider
	serviceProvider          *p.ServiceProvider
	daemonSetProvider        *p.DaemonSetProvider
	persistentVolumeProvider *p.PersistentVolumeProvider
	namespaceProvider        *p.NamespaceProvider
	statefulSetProvider      *p.StatefulSetProvider
	kubeconfig               *string
)

func initWidgets() {
	exampleParagraphWidget = w.NewExampleParagraph()
	nodeListWidget = w.NewNodeListWidget(nodeProvider)
	podListWidget = w.NewKubernetesPods(podProvider)
	summaryWidget = w.NewSummaryWidget(nodeProvider,
		podProvider,
		deploymentProvider,
		serviceProvider,
		daemonSetProvider,
		persistentVolumeProvider,
		namespaceProvider,
		statefulSetProvider,
	)
}

func initObservers() {
	clientSet := newClientSet()
	podProvider = p.NewPodProvider(clientSet)
	deploymentProvider = p.NewDeploymentProvider(clientSet)
	nodeProvider = p.NewNodeProvider(clientSet)
	serviceProvider = p.NewServiceProvider(clientSet)
	daemonSetProvider = p.NewDaemonSetProvider(clientSet)
	persistentVolumeProvider = p.NewPersistentVolumeProvider(clientSet)
	namespaceProvider = p.NewNamespaceProvider(clientSet)
	statefulSetProvider = p.NewStatefulSetProvider(clientSet)
}

func setupGrid() {
	grid = ui.NewGrid()
	grid.Set(
		ui.NewRow(0.25,
			ui.NewCol(0.34, summaryWidget),
			ui.NewCol(0.66, exampleParagraphWidget),
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
	parseFlags()
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

func parseFlags() {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
}

func newClientSet() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	k8s, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return k8s
}

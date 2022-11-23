package main

import (
	"flag"
	"fmt"
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
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	grid                   *ui.Grid
	exampleParagraphWidget w.K8SWidget
	summaryWidget          w.K8SWidget
	nodeListWidget         w.K8SWidget

	updateInterval           = time.Second
	nodeProvider             *p.NodeProvider
	deploymentProvider       *p.DeploymentProvider
	serviceProvider          *p.ServiceProvider
	persistentVolumeProvider *p.PersistentVolumeProvider
	namespaceProvider        *p.NamespaceProvider
	kubeconfig               *string
)

func initWidgets() {
	fmt.Println("initWidgets")
	exampleParagraphWidget = w.NewExampleParagraph()
	nodeListWidget = w.NewNodeListWidget(nodeProvider)
	summaryWidget = w.NewSummaryWidget(nodeProvider,
		deploymentProvider,
		serviceProvider,
		persistentVolumeProvider,
		namespaceProvider,
	)
}

func initObservers(clientSet *kubernetes.Clientset) {
	//podProvider = p.NewPodProvider(clientSet)

	fmt.Println("deploymentProvider")
	deploymentProvider = p.NewDeploymentProvider(clientSet)

	fmt.Println("nodeProvider")
	nodeProvider = p.NewNodeProvider(clientSet)

	fmt.Println("serviceProvider")
	serviceProvider = p.NewServiceProvider(clientSet)

	fmt.Println("persistentVolumeProvider")
	persistentVolumeProvider = p.NewPersistentVolumeProvider(clientSet)

	fmt.Println("namespaceProvider")
	namespaceProvider = p.NewNamespaceProvider(clientSet)

	fmt.Println("Finished Setting up observers")
}

func setupGrid() {
	grid = ui.NewGrid()
	grid.Set(
		ui.NewRow(0.25,
			ui.NewCol(0.34, summaryWidget),
			ui.NewCol(0.66, exampleParagraphWidget),
		),
		ui.NewRow(0.75,
			ui.NewCol(1, nodeListWidget),
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
	fmt.Println("clientSet")
	clientSet := newClientSet()
	fmt.Println(clientSet)

	fmt.Println("observers")
	initObservers(clientSet)

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	fmt.Println("widgets")
	initWidgets()
	fmt.Println("grid")
	setupGrid()
	fmt.Println("term")
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

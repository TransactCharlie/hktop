package main

import (
	ui "github.com/gizak/termui"
	w "github.com/transactcharlie/hktop/src/widgets"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log"
)

var (
	grid *ui.Grid
	exampleParagraphWidget *w.ExampleParagraph
	examplePieWidget *w.ExamplePie
	updateInterval = time.Second
)

func initWidgets() {
	exampleParagraphWidget = w.NewExampleParagraph()
	examplePieWidget = w.NewExamplePie()
}


func setupGrid() {
	grid = ui.NewGrid()
	grid.Set(
		ui.NewRow(1.0/2, exampleParagraphWidget),
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, examplePieWidget),
			ui.NewCol(1.0/2, examplePieWidget),
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
			case "k", "<C-c>":
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
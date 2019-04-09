package widgets

import (
	ui "github.com/gizak/termui/widgets"
)

type ExampleParagraph struct {
	*ui.Paragraph
}

func NewExampleParagraph() *ExampleParagraph {
	p := &ExampleParagraph{Paragraph : ui.NewParagraph()}
	p.Text = "Hello World"
	return p
}

type ExamplePie struct {
	*ui.PieChart
}

func NewExamplePie() *ExamplePie {
	p := &ExamplePie{PieChart: ui.NewPieChart()}
	p.Data = []float64{10, 1, 15}
	return p
}
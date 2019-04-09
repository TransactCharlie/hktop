package widgets

import (
	ui "github.com/gizak/termui/widgets"
)

type ExampleWidget struct {
	*ui.Paragraph
}

func NewExampleWidget() *ExampleWidget {
	p := &ExampleWidget{Paragraph : ui.NewParagraph()}
	p.Text = "Hello World"
	p.Paragraph.Text = "foooooo"
	return p
}
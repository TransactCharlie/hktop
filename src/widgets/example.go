package widgets

import (
	ui "github.com/gizak/termui/v3/widgets"
)

type ExampleParagraph struct {
	*ui.Paragraph
}

func (ep ExampleParagraph) Run() {}
func (ep ExampleParagraph) Stop() bool     {return true}
func (ep ExampleParagraph) Update() error  {return nil}

func NewExampleParagraph() *ExampleParagraph {
	p := &ExampleParagraph{Paragraph : ui.NewParagraph()}
	p.Text = "Hello World"
	return p
}
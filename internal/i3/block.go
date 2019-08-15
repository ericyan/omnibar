package i3

import (
	"barista.run/bar"
	"barista.run/outputs"
	"barista.run/pango"
	"barista.run/pango/icons/mdi"
)

func init() {
	mdi.Load("/usr/local/lib/node_modules/@mdi/font/")
}

// A Block represents a i3bar block.
type Block struct {
	Icon    string
	Text    string
	OnClick func(bar.Event)
}

// Segments implements the bar.Output interface.
func (b *Block) Segments() []*bar.Segment {
	s := outputs.Pango(
		pango.Icon("mdi-"+b.Icon).Large().Rise(-2800),
		pango.Text(" ").Small(),
		b.Text,
	)

	if b.OnClick != nil {
		s.OnClick(b.OnClick)
	}

	return s.Segments()
}

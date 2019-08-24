package keyboard

import (
	"barista.run/bar"
	"barista.run/base/value"
	"barista.run/outputs"

	"github.com/ericyan/omnibar/watchers/ibus"
)

type Module struct {
	formatFunc value.Value // of func(string) bar.Output
}

func New() *Module {
	m := new(Module)
	m.formatFunc.Set(func(in string) bar.Output {
		return outputs.Text(in)
	})
	return m
}

func (m *Module) Stream(s bar.Sink) {
	w, err := ibus.NewWatcher()
	if err != nil {
		s.Error(err)
	}

	updates := w.Watch()

	for u := range updates {
		text := u.Engine
		if u.Label != "" {
			text = u.Label
		}

		format := m.formatFunc.Get().(func(string) bar.Output)
		s.Output(format(text))
	}

	w.Close()
}

func (m *Module) Output(format func(string) bar.Output) *Module {
	m.formatFunc.Set(format)
	return m
}

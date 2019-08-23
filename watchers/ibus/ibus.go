package ibus

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/godbus/dbus"
)

func IBus() (*dbus.Conn, error) {
	dir := os.Getenv("HOME") + "/.config/ibus/bus/"

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(dir + files[0].Name())
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(data), "\n") {
		if kv := strings.SplitN(line, "=", 2); kv[0] == "IBUS_ADDRESS" {
			return dbus.Dial(kv[1])
		}
	}

	return nil, errors.New("unexpected error")
}

type Watcher struct {
	ibus    *dbus.Conn
	state   *State
	signals chan *dbus.Signal
	updates chan State
}

func NewWatcher() (*Watcher, error) {
	conn, err := IBus()
	if err != nil {
		return nil, err
	}

	err = conn.Auth(nil)
	if err != nil {
		return nil, err
	}

	err = conn.Hello()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		conn,
		nil,
		make(chan *dbus.Signal),
		make(chan State),
	}, nil
}

type State struct {
	Engine string
	Label  string
}

func (w *Watcher) Watch() <-chan State {
	w.ibus.Signal(w.signals)
	w.ibus.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.IBus"),
		dbus.WithMatchMember("GlobalEngineChanged"),
	)
	w.ibus.AddMatchSignal(
		dbus.WithMatchInterface("org.freedesktop.IBus.Engine"),
		dbus.WithMatchMember("UpdateProperty"),
	)

	if w.state == nil {
		w.state = &State{"unknown", ""}
	}

	go func() {
		for sig := range w.signals {
			switch sig.Name {
			case "org.freedesktop.IBus.GlobalEngineChanged":
				w.state.Engine = sig.Body[0].(string)
				w.state.Label = ""
			case "org.freedesktop.IBus.Engine.UpdateProperty":
				prop := sig.Body[0].(dbus.Variant).Value().([]interface{})
				if prop[0] != "IBusProperty" {
					continue
				}

				label := prop[4].(dbus.Variant).Value().([]interface{})
				if label[0] != "IBusText" {
					continue
				}

				w.state.Label = label[2].(string)
			}

			w.updates <- *w.state
		}
	}()

	return w.updates
}

func (w *Watcher) Close() {
	w.ibus.RemoveSignal(w.signals)
	close(w.signals)
}

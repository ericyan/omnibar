package main

import (
	"os/exec"

	"barista.run"
	"barista.run/bar"
	"barista.run/modules/clock"
	"barista.run/modules/volume"
	"barista.run/outputs"
	"barista.run/pango"
	"barista.run/pango/icons/mdi"
)

func init() {
	mdi.Load("/usr/local/lib/node_modules/@mdi/font/")
}

func main() {
	barista.Add(volume.DefaultSink().Output(func(vol volume.Volume) bar.Output {
		spacer := pango.Text(" ").Small()

		var seg *bar.Segment
		if vol.Mute {
			seg = outputs.Pango(pango.Icon("mdi-volume-off").Large().Rise(-2800), spacer, "Muted")
		} else {
			seg = outputs.Pango(pango.Icon("mdi-volume-high").Large().Rise(-2800), spacer, pango.Textf("%d%%", vol.Pct()))
		}

		return seg.OnClick(func(e bar.Event) {
			switch e.Button {
			case bar.ButtonLeft:
				vol.SetMuted(!vol.Mute)
			case bar.ButtonRight:
				exec.Command("gnome-control-center", "sound").Run()
			case bar.ScrollUp:
				vol.SetVolume(vol.Vol + (vol.Max-vol.Min)/100)
			case bar.ScrollDown:
				vol.SetVolume(vol.Vol - (vol.Max-vol.Min)/100)
			}
		})
	}))

	barista.Add(clock.Local().OutputFormat("Mon 15:04"))

	barista.Run()
}

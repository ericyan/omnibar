package main

import (
	"os/exec"

	"barista.run"
	"barista.run/bar"
	"barista.run/modules/clock"
	"barista.run/modules/volume"
	"barista.run/outputs"
)

func main() {
	barista.Add(volume.DefaultSink().Output(func(vol volume.Volume) bar.Output {
		var seg *bar.Segment
		if vol.Mute {
			seg = outputs.Text("Vol: Muted")
		} else {
			seg = outputs.Textf("Vol: %d%%", vol.Pct())
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

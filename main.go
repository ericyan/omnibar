package main

import (
	"barista.run"
	"barista.run/bar"
	"barista.run/modules/clock"
	"barista.run/modules/volume"
	"barista.run/outputs"
)

func main() {
	barista.Add(volume.DefaultSink().Output(func(vol volume.Volume) bar.Output {
		if vol.Mute {
			return outputs.Text("Vol: Muted")
		}

		return outputs.Textf("Vol: %d%%", vol.Pct())
	}))

	barista.Add(clock.Local().OutputFormat("Mon 15:04"))

	barista.Run()
}

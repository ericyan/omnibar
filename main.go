package main

import (
	"fmt"
	"os/exec"

	"barista.run"
	"barista.run/bar"
	"barista.run/modules/clock"
	"barista.run/modules/volume"

	"github.com/ericyan/omnibar/internal/i3"
)

func main() {
	barista.Add(volume.DefaultSink().Output(func(vol volume.Volume) bar.Output {
		block := new(i3.Block)

		if vol.Mute {
			block.Icon = "volume-off"
			block.Text = "Muted"
			block.Color = "amber"
		} else {
			block.Icon = "volume-high"
			block.Text = fmt.Sprintf("%d%%", vol.Pct())
		}

		block.OnClick = func(e bar.Event) {
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
		}

		return block
	}))

	barista.Add(clock.Local().OutputFormat("Mon 15:04"))

	barista.Run()
}

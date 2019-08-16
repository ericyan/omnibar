package main

import (
	"fmt"
	"os/exec"

	"barista.run"
	"barista.run/bar"
	"barista.run/modules/clock"
	"barista.run/modules/netinfo"
	"barista.run/modules/volume"
	"barista.run/modules/wlan"

	"github.com/ericyan/omnibar/internal/i3"
)

func main() {
	barista.Add(wlan.Any().Output(func(wifi wlan.Info) bar.Output {
		block := &i3.Block{
			OnClick: func(e bar.Event) {
				switch e.Button {
				case bar.ButtonLeft:
					exec.Command("gnome-control-center", "wifi").Run()
				case bar.ButtonRight:
					exec.Command("gnome-control-center", "network").Run()
				}
			},
		}

		if !wifi.Enabled() {
			block.Icon = "wifi-strength-off"
			block.Text = "Disconnected"
			block.Color = "red"
		}

		if wifi.Connecting() {
			block.Icon = "wifi-strength-outline"
			block.Text = "Connecting..."
			block.Color = "amber"
		}

		if wifi.Connected() {
			block.Icon = "wifi-strength-4"
			block.Text = wifi.SSID
		}

		return block
	}))

	barista.Add(netinfo.New().Output(func(net netinfo.State) bar.Output {
		block := &i3.Block{
			Icon: "ip",
		}

		if len(net.IPs) > 0 {
			block.Text = net.IPs[0].String()
		} else {
			block.Text = "127.0.0.1"
			block.Color = "red"
		}

		return block
	}))

	barista.Add(volume.DefaultSink().Output(func(vol volume.Volume) bar.Output {
		block := &i3.Block{
			Text: fmt.Sprintf("%d%%", vol.Pct()),
		}

		if vol.Mute {
			block.Icon = "volume-off"
			block.Color = "amber"
		} else {
			block.Icon = "volume-high"
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

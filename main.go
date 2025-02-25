package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"barista.run"
	"barista.run/bar"
	"barista.run/modules/battery"
	"barista.run/modules/clock"
	"barista.run/modules/netinfo"
	"barista.run/modules/sysinfo"
	"barista.run/modules/volume"
	"barista.run/modules/wlan"

	"github.com/ericyan/omnibar/internal/i3"
	"github.com/ericyan/omnibar/modules/keyboard"
)

func main() {
	barista.Add(wlan.Any().Output(func(wifi wlan.Info) bar.Output {
		if !wifi.Enabled() {
			return nil
		}

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

		switch {
		case wifi.Connected():
			block.Icon = "wifi-strength-4"
			block.Text = wifi.SSID
		case wifi.Connecting():
			block.Icon = "wifi-strength-outline"
			block.Text = "Connecting..."
			block.Color = "amber"
		default:
			block.Icon = "wifi-strength-off"
			block.Text = "Disconnected"
			block.Color = "red"
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

	barista.Add(sysinfo.New().Output(func(sys sysinfo.Info) bar.Output {
		block := &i3.Block{
			Text: fmt.Sprintf("%.2f", sys.Loads[0]),
		}

		load := (sys.Loads[0] / float64(runtime.NumCPU()))
		switch {
		case load > 1:
			block.Icon = "speedometer"
			block.Color = "red"
		case load > 0.5:
			block.Icon = "speedometer-medium"
			block.Color = "amber"
		default:
			block.Icon = "speedometer-slow"
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

	barista.Add(battery.All().Output(func(bat battery.Info) bar.Output {
		if bat.Status == battery.Disconnected || bat.Status == battery.Unknown {
			return nil
		}

		block := &i3.Block{
			OnClick: func(e bar.Event) {
				if e.Button == bar.ButtonRight {
					exec.Command("gnome-control-center", "power").Run()
				}
			},
		}

		pct := bat.RemainingPct()
		block.Text = fmt.Sprintf("%d%%", pct)

		if bat.PluggedIn() {
			block.Icon = "battery-charging"
			return block
		}

		switch {
		case pct == 100:
			block.Icon = "battery"
		case pct < 10:
			block.Icon = "battery-outline"
		default:
			block.Icon = fmt.Sprintf("battery-%d", pct/10*10)
		}

		switch {
		case pct > 60:
			block.Color = "green"
		case pct > 20:
			block.Color = "amber"
		default:
			block.Color = "red"
		}

		return block
	}))

	barista.Add(keyboard.New().Output(func(text string) bar.Output {
		return &i3.Block{
			Icon: "keyboard",
			Text: text,
			OnClick: func(e bar.Event) {
				if e.Button == bar.ButtonRight {
					exec.Command("gnome-control-center", "region").Run()
				}
			},
		}
	}))

	barista.Add(clock.Local().Output(time.Minute, func(now time.Time) bar.Output {
		return &i3.Block{
			Text: now.Format("Mon 15:04"),
			OnClick: func(e bar.Event) {
				if e.Button == bar.ButtonLeft {
					exec.Command("gsimplecal").Run()
				}
			},
		}
	}))

	barista.Run()
}

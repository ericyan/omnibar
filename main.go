package main

import (
	"barista.run"
	"barista.run/modules/clock"
)

func main() {
	barista.Add(clock.Local().OutputFormat("Mon 15:04"))

	barista.Run()
}

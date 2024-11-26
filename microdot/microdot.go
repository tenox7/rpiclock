// rpiclock driver for Pimoroni Micro Dot pHAT with Lite-On LTP-305
// https://shop.pimoroni.com/en-us/products/microdot-phat
package main

import (
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/jangler/microdotphat-go"
)

type MicroDot struct{}

func (*MicroDot) Open() error {
	slog.Debug("driver: Pimoroni Micro Dot pHAT with Lite-On LTP-305")

	err := microdotphat.Open("")
	if err != nil {
		return err
	}
	microdotphat.SetMirror(true, true)
	return nil
}

func (*MicroDot) Bright() {
	h := time.Now().Local().Hour()
	b := *brNite
	if h > *hrDay && h < *hrNite {
		b = *brDay
	}
	// range 0.0-1.0
	br := float64(b) * 0.01
	microdotphat.SetBrightness(br)
	slog.Debug(fmt.Sprintf("bright: val=%0.2f", br))
}

func (*MicroDot) Close() {
	slog.Debug("closing display")
	microdotphat.Clear()
	microdotphat.Show()
	microdotphat.Close()
}

func (d *MicroDot) DispTime(h, m, s int, pm, syn bool) {
	ind := " "
	sec := " "
	if (s % 2) == 0 {
		sec = ":"
	}
	switch {
	case pm && syn:
		ind = ":"
	case pm:
		ind = "'"
	case syn:
		ind = "."
	}
	microdotphat.WriteString(fmt.Sprintf("%v%02d%v%02d", ind, h, sec, m), 0, 0, false)
	err := microdotphat.Show()
	if err != nil {
		d.Close()
		log.Fatal(err)
	}
}

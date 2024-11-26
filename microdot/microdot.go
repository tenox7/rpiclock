// rpiclock driver for Pimoroni Micro Dot pHAT with Lite-On LTP-305
// https://shop.pimoroni.com/en-us/products/microdot-phat
package main

import (
	"fmt"
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

func (d *MicroDot) Write(s string) error {
	microdotphat.WriteString(s, 0, 0, false)
	return microdotphat.Show()
}

package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/jangler/microdotphat-go"
)

type MicroDot struct{}

func (*MicroDot) Init() error {
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
	microdotphat.SetBrightness(b)
	slog.Debug(fmt.Sprintf("bright: val=%v", b))
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

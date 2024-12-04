// rpiclock dummy display driver, prints out to screen for debug
package main

import (
	"fmt"
	"log/slog"
	"time"
)

type DummyDisp struct{}

func (*DummyDisp) Open() error {
	slog.Debug("driver: dummy / log display")
	return nil
}

func (*DummyDisp) Bright() {
	h := time.Now().Local().Hour()
	b := *brNite
	if h > *hrDay && h < *hrNite {
		b = *brDay
	}
	slog.Debug(fmt.Sprintf("bright: h=%v b=%v", h, b))
}

func (*DummyDisp) Close() {
	slog.Debug("closing display")
}

func (d *DummyDisp) DispTime(h, m, s int, pm, syn bool) {
	slog.Debug(fmt.Sprintf("%02d:%02d:%02d sync=%v pm=%v", h, m, s, syn, pm))
}

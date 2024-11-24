// rpiclock for Adafruit seven segment LED display with backpack (HT16K33)
package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beevik/ntp"
	"github.com/rafalop/sevensegment"
)

var (
	brDay  = flag.Int("br_day", 15, "brightness during day 0..15")
	brNite = flag.Int("br_nite", 0, "brightness during night 0..15")
)

func bright(d *sevensegment.SevenSegment, h int) {
	b := *brNite
	if h > 6 && h < 20 {
		b = *brDay
	}
	d.SetBrightness(b)
}

func tick(d *sevensegment.SevenSegment, l int) {
	h, m, s := time.Now().Local().Clock()
	a := h % 12
	if a == 0 {
		a = 12
	}
	d.SetNum((a * 100) + m)

	var sg [7]bool
	if (s % 2) == 0 {
		sg[sevensegment.IndMidTop] = true
		sg[sevensegment.IndMidBtt] = true
	}
	if h > 11 {
		sg[sevensegment.IndLeftTop] = true
	}
	if l < 3 {
		sg[sevensegment.IndLeftBtt] = true
	}
	d.SetSegments(4, sg)

	if m == 0 && s == 0 {
		bright(d, h)
	}

	d.WriteData()
}

func leap() int {
	r, err := ntp.Query("127.0.0.1")
	if err != nil {
		return 3
	}
	return int(r.Leap)
}

func main() {
	flag.Parse()
	d := sevensegment.NewSevenSegment(0x70)
	bright(d, time.Now().Local().Hour())

	l := leap()
	s := time.NewTicker(time.Second)
	m := time.NewTicker(time.Minute)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.Stop()
		m.Stop()
		d.Clear()
		d.WriteData()
		os.Exit(0)
	}()

	for {
		select {
		case <-m.C:
			l = leap()
		case <-s.C:
			tick(d, l)
		}
	}
}

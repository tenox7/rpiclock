// rpiclock driver for Adafruit seven segment LED display with backpack (HT16K33)
// https://www.adafruit.com/product/1270
package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/rafalop/sevensegment"
)

type SevenSeg struct {
	s *sevensegment.SevenSegment
}

func (s *SevenSeg) Open() error {
	s.s = sevensegment.NewSevenSegment(0x70)
	return nil
}

func (s *SevenSeg) Close() {
	s.s.Clear()
	s.s.WriteData()
}

func (s *SevenSeg) Bright() {
	h := time.Now().Local().Hour()
	b := *brNite
	if h > *hrDay && h < *hrNite {
		b = *brDay
	}
	// range 0-15
	s.s.SetBrightness(int(float64(b) * 0.15))
	slog.Debug(fmt.Sprintf("bright: val=%v", b))
}

func (ss *SevenSeg) DispTime(h, m, s int, pm, syn bool) {
	ss.s.SetNum((h * 100) + m)

	var sg [7]bool
	if (s % 2) == 0 {
		sg[sevensegment.IndMidTop] = true
		sg[sevensegment.IndMidBtt] = true
	}
	if pm {
		sg[sevensegment.IndLeftTop] = true
	}
	if syn {
		sg[sevensegment.IndLeftBtt] = true
	}
	ss.s.SetSegments(4, sg)
	ss.s.WriteData()
}

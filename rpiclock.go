package main

import (
	"time"

	"github.com/rafalop/sevensegment"
)

func bright(d *sevensegment.SevenSegment, h int) {
	b := 0
	if h > 7 && h < 20 {
		b = 15
	}
	d.SetBrightness(b)
}

func tick(d *sevensegment.SevenSegment) {
	h, m, s := time.Now().Local().Clock()
	d.SetNum(((h % 12) * 100) + m)

	// segments: 0,1=tick 2=tl 3=bl 4=tr
	var sg [7]bool
	if (s % 2) == 0 {
		sg[0] = true
		sg[1] = true
	}
	if h > 11 {
		sg[2] = true
	}
	d.SetSegments(4, sg)

	if m == 0 && s == 0 {
		bright(d, h)
	}

	d.WriteData()
}

func main() {
	d := sevensegment.NewSevenSegment(0x70)
	bright(d, time.Now().Local().Hour())

	s := time.NewTicker(time.Second)

	for {
		select {
		case <-s.C:
			tick(d)
		}
	}
}

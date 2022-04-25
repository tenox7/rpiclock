/*
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
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

	// segments: 0,1=tick 2=topLeft 3=bottomLeft 4=topRight
	var sg [7]bool
	if (s % 2) == 0 {
		sg[0] = true
		sg[1] = true
	}
	if h > 11 {
		sg[2] = true
	}
	if l < 3 {
		sg[3] = true
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
	s := time.NewTicker(time.Second)
	n := make(chan int)
	l := 0

	bright(d, time.Now().Local().Hour())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		d.Clear()
		d.WriteData()
		os.Exit(0)
	}()

	go func(c chan<- int) {
		for {
			c <- leap()
			time.Sleep(60 * time.Second)
		}
	}(n)

	for {
		select {
		case l = <-n:
		case <-s.C:
			tick(d, l)
		}
	}
}

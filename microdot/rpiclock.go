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

// rpiclock for Pimoroni Micro Dot pHAT with Lite-On LTP-305
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beevik/ntp"
	"github.com/jangler/microdotphat-go"
)

var (
	brDay  = flag.Float64("br_day", 1.0, "brightness during day 0.0-1.0")
	brNite = flag.Float64("br_nite", 0.3, "brightness during night 0.0-1.0")
)

func bright(h int) {
	b := *brNite
	if h > 6 && h < 20 {
		b = *brDay
	}
	microdotphat.SetBrightness(b)
}

func tick(l int) {
	h, m, s := time.Now().Local().Clock()
	a := h % 12
	if a == 0 {
		a = 12
	}
	ind := " "
	sec := " "
	if (s % 2) == 0 {
		sec = ":"
	}
	switch {
	case h > 11 && l < 3:
		ind = ":"
	case h > 11:
		ind = "'"
	case l < 3:
		ind = "."
	}
	if m == 0 && s == 0 {
		bright(h)
	}
	microdotphat.WriteString(fmt.Sprintf("%v%02d%v%02d", ind, a, sec, m), 0, 0, false)
	err := microdotphat.Show()
	if err != nil {
		log.Println(err)
		clear()
	}
}

func clear() {
	microdotphat.Clear()
	microdotphat.Show()
	time.Sleep(time.Millisecond)
	microdotphat.Close()
	os.Exit(0)
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

	err := microdotphat.Open("")
	if err != nil {
		log.Fatal(err)
	}
	microdotphat.SetMirror(true, true)
	bright(time.Now().Local().Hour())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		clear()
	}()

	n := make(chan int)
	go func(c chan<- int) {
		for {
			c <- leap()
			time.Sleep(60 * time.Second)
		}
	}(n)

	l := 0
	s := time.NewTicker(time.Second)
	for {
		select {
		case l = <-n:
		case <-s.C:
			tick(l)
		}
	}
}

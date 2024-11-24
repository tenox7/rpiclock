// rpiclock for Pimoroni Micro Dot pHAT with Lite-On LTP-305
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/beevik/ntp"
	"github.com/jangler/microdotphat-go"
)

var (
	brDay  = flag.Float64("br_day", 1.0, "brightness during day 0.0-1.0")
	brNite = flag.Float64("br_nite", 0.3, "brightness during night 0.0-1.0")
)

type RPIClock struct {
	l int
	sync.Mutex
}

func bright(h int) {
	b := *brNite
	if h > 6 && h < 20 {
		b = *brDay
	}
	microdotphat.SetBrightness(b)
}

func (r *RPIClock) tick() {
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
	r.Lock()
	l := r.l
	r.Unlock()
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

	r := RPIClock{}
	r.leap()
	s := time.NewTicker(time.Second)
	m := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-m.C:
			go r.leap()
		case <-s.C:
			r.tick()
		}
	}
}

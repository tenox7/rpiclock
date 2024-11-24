// rpiclock for Pimoroni Micro Dot pHAT with Lite-On LTP-305
package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/beevik/ntp"
	"github.com/jangler/microdotphat-go"
)

// TODO
// - 24 hours as flag

var (
	brDay  = flag.Float64("br_day", 1.0, "brightness during day 0.0-1.0")
	brNite = flag.Float64("br_nite", 0.3, "brightness during night 0.0-1.0")
	hrDay  = flag.Int("hr_day", 6, "bright display / day start hour (24h)")
	hrNite = flag.Int("hr_nite", 20, "dim display / nite start hour (24h)")
	ntpq   = flag.Duration("ntpq", time.Minute, "ntp sync status query interval")
	debug  = flag.Bool("debug", false, "debug logging")
)

type RPIClock struct {
	synchronized bool
	sync.Mutex
}

func (_ *RPIClock) bright() {
	h := time.Now().Local().Hour()
	b := *brNite
	if h > *hrDay && h < *hrNite {
		b = *brDay
	}
	microdotphat.SetBrightness(b)
	slog.Debug(fmt.Sprintf("bright: val=%v", b))
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
	syn := r.synchronized
	r.Unlock()
	switch {
	case h > 11 && syn:
		ind = ":"
	case h > 11:
		ind = "'"
	case syn:
		ind = "."
	}
	microdotphat.WriteString(fmt.Sprintf("%v%02d%v%02d", ind, a, sec, m), 0, 0, false)
	err := microdotphat.Show()
	if err != nil {
		slog.Error(err.Error())
		r.clear()
	}
}

func (_ *RPIClock) clear() {
	slog.Debug("clearing display")
	microdotphat.Clear()
	microdotphat.Show()
	time.Sleep(time.Millisecond)
	microdotphat.Close()
	os.Exit(0)
}

func (r *RPIClock) ntpq() {
	n, err := ntp.Query("127.0.0.1")
	r.Mutex.Lock()
	defer func() { slog.Debug(fmt.Sprintf("ntp: sync=%v err=%v", r.synchronized, err)); r.Mutex.Unlock() }()
	if err != nil || n.Leap > 2 {
		r.synchronized = false
		return
	}
	r.synchronized = true
}

func main() {
	flag.Parse()
	if *debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	slog.Debug("RPI Clock Starting Up")
	slog.Debug("Using Pimoroni Micro Dot pHAT with Lite-On LTP-305")

	err := microdotphat.Open("")
	if err != nil {
		log.Fatal(err)
	}
	microdotphat.SetMirror(true, true)

	r := RPIClock{}
	r.ntpq()
	r.bright()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		r.clear()
	}()

	tT := time.NewTicker(time.Second)
	nT := time.NewTicker(*ntpq)
	bT := time.NewTicker(time.Hour)
	for {
		select {
		case <-tT.C:
			r.tick()
		case <-nT.C:
			go r.ntpq()
		case <-bT.C:
			r.bright()
		}
	}
}

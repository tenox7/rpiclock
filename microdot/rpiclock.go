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
)

var (
	t24h   = flag.Bool("t24h", false, "use 24h time format")
	brDay  = flag.Float64("br_day", 1.0, "brightness during day 0.0-1.0")
	brNite = flag.Float64("br_nite", 0.3, "brightness during night 0.0-1.0")
	hrDay  = flag.Int("hr_day", 6, "bright display / day start hour (24h)")
	hrNite = flag.Int("hr_nite", 20, "dim display / nite start hour (24h)")
	ntpq   = flag.Duration("ntpq", time.Minute, "ntp sync status query interval")
	dspDrv = flag.String("disp", "microdot", "display driver: microdot")
	debug  = flag.Bool("debug", false, "debug logging")
)

type DisplayDriver interface {
	Init() error
	Close()
	Bright()
	Write(string) error
}

type RPIClock struct {
	disp    DisplayDriver
	ntpSync bool
	sync.Mutex
}

func (r *RPIClock) tick() {
	h, m, s := time.Now().Local().Clock()
	a := h
	if !*t24h {
		a = h % 12
		if a == 0 {
			a = 12
		}
	}
	ind := " "
	sec := " "
	if (s % 2) == 0 {
		sec = ":"
	}
	r.Lock()
	syn := r.ntpSync
	r.Unlock()
	switch {
	case h > 11 && syn:
		ind = ":"
	case h > 11:
		ind = "'"
	case syn:
		ind = "."
	}
	err := r.disp.Write(fmt.Sprintf("%v%02d%v%02d", ind, a, sec, m))
	if err != nil {
		r.disp.Close()
		log.Fatal(err)
	}
}

func (r *RPIClock) ntpCheck() {
	n, err := ntp.Query("127.0.0.1")
	r.Mutex.Lock()
	defer func() {
		slog.Debug(fmt.Sprintf("ntp: sync=%v err=%v", r.ntpSync, err))
		r.Mutex.Unlock()
	}()
	if err != nil || n.Leap > 2 {
		r.ntpSync = false
		return
	}
	r.ntpSync = true
}

func main() {
	flag.Parse()
	if *debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	slog.Debug("RPI Clock Starting Up")

	r := RPIClock{}
	r.ntpCheck()

	switch *dspDrv {
	case "microdot":
		r.disp = &MicroDot{}
	default:
		log.Fatalf("unsupported display driver: %v", *dspDrv)
	}
	err := r.disp.Init()
	if err != nil {
		log.Fatalf("Unable to initialize %v: %v", *dspDrv, err)
	}
	r.disp.Bright()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		r.disp.Close()
		os.Exit(0)
	}()

	tT := time.NewTicker(time.Second)
	nT := time.NewTicker(*ntpq)
	bT := time.NewTicker(time.Hour)
	for {
		select {
		case <-tT.C:
			r.tick()
		case <-nT.C:
			go r.ntpCheck()
		case <-bT.C:
			r.disp.Bright()
		}
	}
}

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
	t24h   = flag.Bool("24h", false, "use 24h time format")
	brDay  = flag.Int("br_day", 100, "brightness during day 0-100")
	brNite = flag.Int("br_nite", 30, "brightness during night 0-100")
	hrDay  = flag.Int("hr_day", 6, "bright display / day start hour (24h)")
	hrNite = flag.Int("hr_nite", 20, "dim display / nite start hour (24h)")
	ntpq   = flag.Duration("ntpq", time.Minute, "ntp sync status query interval")
	dspDrv = flag.String("disp", "sevensegment", "display driver: sevensegment|microdot")
	debug  = flag.Bool("debug", false, "debug logging")
)

type DisplayDriver interface {
	Open() error
	Close()
	Bright()
	DispTime(h, m, s int, pm, syn bool)
}

type RPIClock struct {
	disp       DisplayDriver
	ntpIsSynch bool
	sync.Mutex
}

func (r *RPIClock) ntpStat() bool {
	r.Lock()
	defer r.Unlock()
	return r.ntpIsSynch
}

func (r *RPIClock) ntpCheck() {
	n, err := ntp.Query("127.0.0.1")
	r.Mutex.Lock()
	defer func() {
		slog.Debug(fmt.Sprintf("ntp: sync=%v err=%v", r.ntpIsSynch, err))
		r.Mutex.Unlock()
	}()
	if err != nil || n.Leap > 2 {
		r.ntpIsSynch = false
		return
	}
	r.ntpIsSynch = true
}

func (r *RPIClock) tick() {
	h, m, s := time.Now().Local().Clock()
	pm := false
	if h > 11 {
		pm = true
	}
	if !*t24h {
		h = h % 12
		if h == 0 {
			h = 12
		}
	}
	r.disp.DispTime(h, m, s, pm, r.ntpStat())
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
	case "sevensegment":
		r.disp = &SevenSeg{}
	case "microdot":
		r.disp = &MicroDot{}
	default:
		log.Fatalf("unsupported display driver: %v", *dspDrv)
	}
	err := r.disp.Open()
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

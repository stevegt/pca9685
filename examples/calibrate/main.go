package main

import (
	"strconv"
	"strings"
	"time"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"

	"github.com/chzyer/readline"
	. "github.com/stevegt/goadapt"
	"github.com/stevegt/pca9685"
)

func main() {

	// this is pretty standard gobot "metal" style boilerplate
	rpi := raspi.NewAdaptor()
	pca := i2c.NewPCA9685Driver(rpi, i2c.WithBus(1), i2c.WithAddress(0x40))

	// initialize the chip
	err := pca.Start()
	Ck(err)

	// create a Servo object on channel 15
	s := pca9685.NewServo(pca, 15)

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "n {min pulse uS} | x {max pulse uS} | a {angle degrees}: ",
		HistoryFile:     ".calibrate-history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	Ck(err)
	defer l.Close()
	l.CaptureExitSignal()

	var a float32
	for {
		line, err := l.Readline()
		Ck(err)
		parts := strings.Split(line, " ")
		if len(parts) < 1 {
			continue
		}
		switch parts[0] {
		case "n":
			us64, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				continue
			}
			s.PulseMin = time.Duration(us64) * time.Microsecond
			Pf("pulse min %v", s.PulseMin)
		case "x":
			us64, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				continue
			}
			s.PulseMax = time.Duration(us64) * time.Microsecond
			Pf("pulse max %v", s.PulseMax)
		case "a":
			a64, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				continue
			}
			a = float32(a64)
		}
		t := s.Angle(a)
		Pf("%v\n", t)
	}

	/*
		for i := 100; i < 3000; i += 10 {
			s.PulseMin = time.Duration(i) * time.Microsecond
			s.PulseMax = s.PulseMin + 100*time.Microsecond
			t := s.Angle(s.AngleMin)
			Pf("%v\n", t)
			time.Sleep(100 * time.Millisecond)
			t = s.Angle(s.AngleMax)
			Pf("%v\n", t)
			time.Sleep(100 * time.Millisecond)
		}
	*/

	/*
		var ctr float32 = 45.0
		var jitter float32 = 2.0
		s.Angle(ctr)
		// for ms := s.PulseMin; a < s.AngleMax; a++ {
		for {
			// t := s.Angle(a)
			// Pf("%v\n", t)

			// jitter servo +/- 1 degree
			t := s.Angle(ctr - jitter)
			Pf("%v\n", t)
			time.Sleep(100 * time.Millisecond)
			t = s.Angle(ctr + jitter)
			Pf("%v\n", t)
			time.Sleep(100 * time.Millisecond)
		}
	*/
}

package main

import (
	"bufio"
	"log"
	"os"
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
		Prompt:          "w (calibration wizard) | n {min pulse uS} | x {max pulse uS} | a {angle degrees}: ",
		HistoryFile:     ".calibrate-history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	Ck(err)
	defer l.Close()
	// l.CaptureExitSignal()

	var a float32
	for {
		line, err := l.Readline()
		Ck(err)
		parts := strings.Split(line, " ")
		if len(parts) < 1 {
			continue
		}
		switch parts[0] {
		case "w":
			wizard(s)
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
}

func wizard(s *pca9685.Servo) {
	// find min pulse
	jiggle(s, 100, 2000, 100)
	time.Sleep(3 * time.Second)
	jiggle(s, int(s.PulseMin/time.Microsecond), int(s.PulseMax/time.Microsecond), 10)
	pulseMin := s.PulseMax

	// find max pulse
	jiggle(s, 3000, 1000, -100)
	time.Sleep(3 * time.Second)
	jiggle(s, int(s.PulseMax/time.Microsecond), int(s.PulseMin/time.Microsecond), -10)
	pulseMax := s.PulseMin

	s.PulseMin = pulseMin
	s.PulseMax = pulseMax

	// measure angle
	s.Angle(0)
	angleMin := inputFloat32("raw protractor angle")
	s.Angle(s.AngleRange)
	angleMax := inputFloat32("raw protractor angle")
	s.AngleRange = angleMax - angleMin

	// show results
	Pf("pulse min %v max %v angle range %.0f degrees\n", s.PulseMin, s.PulseMax, s.AngleRange)
	for i := 0; i < 3; i++ {
		s.Angle(0)
		time.Sleep(1 * time.Second)
		s.Angle(s.AngleRange)
		time.Sleep(1 * time.Second)
	}

}

func hitkey() bool {
	var buf []byte
	n, err := os.Stdin.Read(buf)
	Ck(err)
	if n > 0 {
		return true
	}
	return false
}

func jiggle(s *pca9685.Servo, start, stop, step int) {
	wait := 1000 * time.Millisecond
	// os.Stdin.SetReadDeadline(wait)
	input := make(chan rune, 1)
	quit := make(chan bool, 1)
	go readKey(input, quit)
	Pl("hit enter when movement starts")
loop:
	for i := start; ; i += step {
		if step > 0 && i > stop {
			break
		}
		if step < 0 && i < stop {
			break
		}
		s.PulseMin = time.Duration(i) * time.Microsecond
		s.PulseMax = time.Duration(i+100) * time.Microsecond
		s.Angle(0)
		time.Sleep(wait)
		s.Angle(90)
		select {
		case <-input:
			Pl("break")
			break loop
		case <-time.After(wait):
			continue
		}
	}
	quit <- true
}

func inputFloat32(prompt string) (res float32) {
	Pf("%s:\n", prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	Ck(err)
	str := strings.TrimSpace(line)
	res64, err := strconv.ParseFloat(str, 32)
	Ck(err)
	return float32(res64)
}

func readKey(input chan rune, quit chan bool) {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	select {
	case input <- char:
		return
	case <-quit:
		return
	}
}

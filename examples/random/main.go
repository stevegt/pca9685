package main

import (
	"math/rand"
	"time"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"

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

	for {
		// move servo around randomly
		a := rand.Float32() * s.AngleRange
		t := s.Angle(a)

		// print trace info
		Pf("%v\n", t)

		time.Sleep(time.Second)
	}
}

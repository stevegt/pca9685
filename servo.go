package pca9685

import (
	"time"

	"gobot.io/x/gobot/drivers/i2c"
)

type Servo struct {
	driver   *i2c.PCA9685Driver
	channel  int
	Freq     float32
	PulseMin time.Duration
	PulseMax time.Duration
	AngleMin float32
	AngleMax float32
}

func NewServo(pca9685 *i2c.PCA9685Driver, channel int) (s *Servo) {
	s = &Servo{
		driver:   pca9685,
		channel:  channel,
		Freq:     50,
		PulseMin: 1 * time.Millisecond,
		PulseMax: 2 * time.Millisecond,
		AngleMin: 0,
		AngleMax: 90,
	}
	return
}

type Trace struct {
	Period     time.Duration
	Angle      float32
	AnglePct   float32
	PulseRange time.Duration
	OffTime    time.Duration
	OffTicks   uint16
}

func (s *Servo) Angle(angle float32) (t Trace) {
	t.Angle = angle
	t.Period = time.Duration(float32(time.Second) / s.Freq)
	t.AnglePct = (angle - s.AngleMin) / s.AngleMax
	t.PulseRange = s.PulseMax - s.PulseMin
	t.OffTime = s.PulseMin + time.Duration(float32(t.PulseRange)*t.AnglePct)
	t.OffTicks = uint16(4096 * t.OffTime / t.Period)
	s.driver.SetPWM(s.channel, 0, t.OffTicks)
	return
}

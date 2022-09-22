# pca9685

Yet another library for driving servos via a PCA9685 chip.  These chips are used in many popular i2c-interface PWM controller boards combatible with the Arduino, Raspberry Pi, and other single board computer and microcontroller boards.  Clean and simple, uses gobot for backend.  Includes a simple servo calibration tool.  See ./examples.

Written in part because the gobot ServoDriver library only uses 8 bits
of resolution, while the pca9685 chip supports 12 bits. This library
uses all 12 bits.

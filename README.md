# pca9685

Yet another library for driving servos via a PCA9685.  Clean and
simple, uses gobot for backend.  Includes a simple servo calibration
tool.  See ./examples.

Written in part because the gobot ServoDriver library only uses 8 bits
of resolution, while the pca9685 chip supports 12 bits. This library
uses all 12 bits.

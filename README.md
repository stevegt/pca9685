# pca9685

Yet another library for driving servos via a PCA9685.  Clean and
simple, uses gobot for backend.  See ./examples.

Written in part because the gobot ServoDriver library only has uses 8
bits of resolution, while the pca9685 supports 12 bits.

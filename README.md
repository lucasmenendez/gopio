# gopio
Simple golang GPIO library. Abstraction of [GPIO Sysfs Interface](https://www.kernel.org/doc/Documentation/gpio/sysfs.txt).

## API Reference

#### gopio.NewPin(pin int)
Return gpio pin selected struct.

#### pin.Close()
Unexport gpio.

#### pin.SetMode(mode string)
Set pin mode between "in" or "out".

#### pin.GetMode()
Return pin mode.

#### pin.SetValue(value int)
Set pin value between 0 or 1.

#### pin.GetValue()
Return pin value.

// gopio package implements simple API in go to work with sysfs gpios. Read more
// about sysfs here: https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
package gopio

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Pin struct {
	id string
}

const gpioPath string = "/sys/class/gpio"

// New function initializes gpio by pin id
func New(id int) (pin *Pin, e error) {
	pin.id = strconv.Itoa(id)

	return pin, pin.init()
}

// init function check if current pin id exists. If not, create socket fd
// and export the pin.
func (pin *Pin) init() (e error) {
	var (
		p  string = fmt.Sprintf("%s/gpio%s", gpioPath, pin.id)
		ex string = fmt.Sprintf("%s/export", gpioPath)
	)

	if _, e = os.Stat(p); os.IsNotExist(e) {
		var fd *os.File
		if fd, e = os.OpenFile(ex, os.O_APPEND|os.O_WRONLY, 0777); e != nil {
			return e
		}
		defer fd.Close()

		if _, e = fd.WriteString(pin.id); e != nil {
			return e
		}
	}

	return e
}

// Close function unexports pin socket.
func (pin *Pin) Close() (e error) {
	var uex string = fmt.Sprintf("%s/unexport", gpioPath)

	var fd *os.File
	if fd, e = os.OpenFile(uex, os.O_APPEND|os.O_WRONLY, 0777); e != nil {
		return e
	}
	defer fd.Close()

	if _, e = fd.WriteString(pin.id); e != nil {
		return e
	}
	return nil
}

// SetMode receives a string that contains new mode to set current pin. First
// check if received mode is valid and then write into the socket.
func (pin *Pin) SetMode(m string) (e error) {
	if m != "in" && m != "out" {
		return errors.New("Invalid mode.")
	}

	var (
		p  string = fmt.Sprintf("%s/gpio%s", gpioPath, pin.id)
		dp string = fmt.Sprintf("%s/direction", p)
	)

	if _, e = os.Stat(p); e != nil {
		return e
	}

	var fd *os.File
	if fd, e = os.OpenFile(dp, os.O_APPEND|os.O_WRONLY, 0777); e != nil {
		return e
	}
	defer fd.Close()

	if _, e = fd.WriteString(m); e != nil {
		return e
	}

	return nil
}

// GetMode returns current pin mode reading it from direction descriptor
// associated to the pin.
func (pin *Pin) GetMode() (m string, e error) {
	var (
		p  string = fmt.Sprintf("%s/gpio%s", gpioPath, pin.id)
		dp string = fmt.Sprintf("%s/direction", p)
	)

	if _, e = os.Stat(p); e != nil {
		return m, e
	}

	var fd *os.File
	if fd, e = os.OpenFile(dp, os.O_RDONLY, 0777); e != nil {
		return m, e
	}
	defer fd.Close()

	var c int
	data := make([]byte, 3)
	if c, e = fd.Read(data); e != nil {
		return m, e
	}

	return strings.TrimSpace(string(data[:c])), e
}

// SetValue receives an int that contains new value to set current pin. First
// check if received value is valid and then write into the socket.
func (pin *Pin) SetValue(v int) (e error) {
	if v != 1 && v != 0 {
		return errors.New("Invalid value.")
	}

	var (
		p  string = fmt.Sprintf("%s/gpio%s", gpioPath, pin.id)
		vp string = fmt.Sprintf("%s/value", p)
	)

	if _, e = os.Stat(p); e != nil {
		return e
	}

	var fd *os.File
	if fd, e = os.OpenFile(vp, os.O_APPEND|os.O_WRONLY, 0777); e != nil {
		return e
	}
	defer fd.Close()

	if _, e = fd.WriteString(strconv.Itoa(v)); e != nil {
		return e
	}

	return nil
}

// GetValue returns current pin value reading it from direction descriptor
// associated to the pin.
func (pin *Pin) GetValue() (v int, e error) {
	var (
		p  string = fmt.Sprintf("%s/gpio%s", gpioPath, pin.id)
		vp string = fmt.Sprintf("%s/value", p)
	)

	if _, e = os.Stat(p); e != nil {
		return -1, e
	}

	var fd *os.File
	if fd, e = os.OpenFile(vp, os.O_RDONLY, 0777); e != nil {
		return -1, e
	}
	defer fd.Close()

	d := make([]byte, 1)
	if _, e = fd.Read(d); e != nil {
		return -1, e
	}

	var res string = strings.TrimSpace(string(d))
	if v, e = strconv.Atoi(res); e != nil {
		return -1, e
	}

	return v, nil
}

package gopio

import (
	"os"
	"strconv"
	"strings"
)

type Pin struct {
	id string
}

const gpioPath = "/sys/class/gpio/"

func NewPin(id int) *Pin {
	pin := new(Pin)
	pin.id = strconv.Itoa(id)

	pin.Init()
	return pin
}

func (pin *Pin) Init() {
	var err error

	if _, err = os.Stat(gpioPath + "gpio" + pin.id); os.IsNotExist(err) {
		var file *os.File

		if file, err = os.OpenFile(gpioPath+"export", os.O_APPEND|os.O_WRONLY, 0777); err != nil {
			panic(err)
		}
		defer file.Close()

		if _, err = file.WriteString(pin.id); err != nil {
			panic(err)
		}
	}
	return
}

func (pin *Pin) SetMode(mode string) {
	if mode != "in" && mode != "out" {
		panic("Invalid mode.")
	}

	var err error
	if _, err = os.Stat(gpioPath + "gpio" + pin.id); err != nil {
		panic(err)
	}

	var file *os.File
	if file, err = os.OpenFile(gpioPath+"gpio"+pin.id+"/direction", os.O_APPEND|os.O_WRONLY, 0777); err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = file.WriteString(mode); err != nil {
		panic(err)
	}

	return
}

func (pin *Pin) GetMode() string {
	var err error
	if _, err = os.Stat(gpioPath + "gpio" + pin.id); err != nil {
		panic(err)
	}

	var file *os.File
	if file, err = os.OpenFile(gpioPath+"gpio"+pin.id+"/direction", os.O_RDONLY, 0777); err != nil {
		panic(err)
	}
	defer file.Close()

	var count int
	data := make([]byte, 3)
	if count, err = file.Read(data); err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(data[:count]))
}

func (pin *Pin) SetValue(value int) {
	if value != 1 && value != 0 {
		panic("Invalid value.")
	}

	var err error
	if _, err = os.Stat(gpioPath + "gpio" + pin.id); err != nil {
		panic(err)
	}

	var file *os.File
	if file, err = os.OpenFile(gpioPath+"gpio"+pin.id+"/value", os.O_APPEND|os.O_WRONLY, 0777); err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = file.WriteString(strconv.Itoa(value)); err != nil {
		panic(err)
	}

	return
}

func (pin *Pin) GetValue() int {
	var err error
	if _, err = os.Stat(gpioPath + "gpio" + pin.id); err != nil {
		panic(err)
	}

	var file *os.File
	if file, err = os.OpenFile(gpioPath+"gpio"+pin.id+"/value", os.O_RDONLY, 0777); err != nil {
		panic(err)
	}
	defer file.Close()

	data := make([]byte, 1)
	if _, err = file.Read(data); err != nil {
		panic(err)
	}

	var res int
	if res, err = strconv.Atoi(strings.TrimSpace(string(data))); err != nil {
		panic(err)
	}

	return res
}

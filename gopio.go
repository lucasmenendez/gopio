package gopio

import (
	"os"
	"strconv"
	"strings"
	"errors"
)

type Pin struct {
	id string
}

const gpioPath = "/sys/class/gpio/"

func NewPin(id int) (error, *Pin) {
	pin := new(Pin)
	pin.id = strconv.Itoa(id)

	err := pin.Init()
	return err, pin
}

func (pin *Pin) Init() error {
	var err error

	if _, err = os.Stat(gpioPath + "gpio" + pin.id); os.IsNotExist(err) {
		var file *os.File

		if file, err = os.OpenFile(gpioPath+"export", os.O_APPEND|os.O_WRONLY, 0777); err != nil {
			return err
		}
		defer file.Close()

		if _, err = file.WriteString(pin.id); err != nil {
			return err
		}
	}
	return nil
}

func (pin *Pin) Close() error {
	var err error
	var file *os.File

	if file, err = os.OpenFile(gpioPath+"unexport", os.O_APPEND|os.O_WRONLY, 0777); err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(pin.id); err != nil {
		return err
	}
	return nil
}

func (pin *Pin) SetMode(mode string) error {
	if mode != "in" && mode != "out" {
		return errors.New("Invalid mode.")
	}

	var err error
	if _, err = os.Stat(gpioPath + "gpio" + pin.id); err != nil {
		return err
	}

	var file *os.File
	if file, err = os.OpenFile(gpioPath+"gpio"+pin.id+"/direction", os.O_APPEND|os.O_WRONLY, 0777); err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(mode); err != nil {
		return err
	}

	return nil
}

func (pin *Pin) GetMode() (error, string) {
	var err error
	if _, err = os.Stat(gpioPath + "gpio" + pin.id); err != nil {
		return err, ""
	}

	var file *os.File
	if file, err = os.OpenFile(gpioPath+"gpio"+pin.id+"/direction", os.O_RDONLY, 0777); err != nil {
		return err, ""
	}
	defer file.Close()

	var count int
	data := make([]byte, 3)
	if count, err = file.Read(data); err != nil {
		return err, ""
	}

	return nil, strings.TrimSpace(string(data[:count]))
}

func (pin *Pin) SetValue(value int) error {
	if value != 1 && value != 0 {
		return errors.New("Invalid value.")
	}

	var err error
	if _, err = os.Stat(gpioPath + "gpio" + pin.id); err != nil {
		return err
	}

	var file *os.File
	if file, err = os.OpenFile(gpioPath+"gpio"+pin.id+"/value", os.O_APPEND|os.O_WRONLY, 0777); err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(strconv.Itoa(value)); err != nil {
		return err
	}

	return nil
}

func (pin *Pin) GetValue() (error, int) {
	var err error
	
	if _, err = os.Stat(gpioPath + "gpio" + pin.id); err != nil {
		return err, -1
	}

	var file *os.File
	if file, err = os.OpenFile(gpioPath+"gpio"+pin.id+"/value", os.O_RDONLY, 0777); err != nil {
		return err, -1
	}
	defer file.Close()

	data := make([]byte, 1)
	if _, err = file.Read(data); err != nil {
		return err, -1
	}

	var res int
	if res, err = strconv.Atoi(strings.TrimSpace(string(data))); err != nil {
		return err, -1
	}

	return nil, res
}

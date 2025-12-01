package day_01

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	pEnd   = 0x454E44
	pClick = 0x434C49434B
)

type Dial struct {
	logger         *log.Logger
	position       int
	size           int
	password       int
	passwordMethod int
}

func NewDial(passwordMethod string) (*Dial, error) {
	dial := &Dial{
		logger:         log.New(os.Stdout, "", log.Lshortfile),
		position:       50,
		size:           100,
		password:       0,
		passwordMethod: 0,
	}
	switch strings.ToUpper(passwordMethod) {
	case "END":
		dial.passwordMethod = pEnd
	case "CLICK":
		dial.passwordMethod = pClick
	default:
		return &Dial{}, fmt.Errorf("unhandled password method %s", passwordMethod)
	}
	return dial, nil
}

func (d *Dial) GetPassword(reader io.Reader) (int, error) {
	d.logger.Printf("The dial starts at position %d\n", d.position)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		err := d.turnDialUsingInstruction(line)
		if err != nil {
			return 0, err
		}
	}
	d.logger.Printf("The password is %d\n", d.password)
	return d.password, nil
}

func (d *Dial) turnDialUsingInstruction(line string) error {
	sign := 1
	direction := line[:1]
	if direction != "L" && direction != "R" {
		return fmt.Errorf("unknown instruction %s", direction)
	}
	if direction == "L" {
		sign = -1
	}
	amount, err := strconv.Atoi(line[1:])
	if err != nil {
		return err
	}
	for i := range amount {
		next := d.position + sign
		if next == d.size && sign == 1 {
			next = 0
		} else if next == -1 && sign == -1 {
			next = d.size - 1
		} else {
			// do nothing
		}
		d.position = next
		if i+1 == amount && d.position == 0 {
			d.password += 1
		}
	}
	d.logger.Printf("The dial is rotated %s to point at %d\n", line, d.position)
	return nil
}

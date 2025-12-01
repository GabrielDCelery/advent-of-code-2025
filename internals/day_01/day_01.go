package day_01

import (
	"bufio"
	"container/ring"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const defaultRingSize = 100
const defaultRingPosition = 50

const (
	passwordMethodEnd   = 0x454E44
	passwordMethodClick = 0x434C49434B
)

type Dial struct {
	logger         *log.Logger
	ring           *ring.Ring
	password       int
	passwordMethod int
}

func NewDial(passwordMethod string) (*Dial, error) {
	ringSize := defaultRingSize
	ringPosition := defaultRingPosition
	dial := &Dial{
		logger:         log.New(os.Stdout, "", log.Lshortfile),
		ring:           ring.New(ringSize),
		password:       0,
		passwordMethod: 0,
	}
	for i := range ringSize {
		dial.ring.Value = i
		dial.ring = dial.ring.Next()
	}
	dial.ring = dial.ring.Move(ringPosition)
	switch strings.ToUpper(passwordMethod) {
	case "END":
		dial.passwordMethod = passwordMethodEnd
	case "CLICK":
		dial.passwordMethod = passwordMethodClick
	default:
		return &Dial{}, fmt.Errorf("unhandled password method %s", passwordMethod)
	}
	return dial, nil
}

func (d *Dial) GetPassword(reader io.Reader) (int, error) {
	d.logger.Printf("The dial starts at position %d\n", d.ring.Value.(int))
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
	dirRight := true
	direction := line[:1]
	if direction != "L" && direction != "R" {
		return fmt.Errorf("unknown instruction %s", direction)
	}
	if direction == "L" {
		dirRight = false
	}
	amount, err := strconv.Atoi(line[1:])
	if err != nil {
		return err
	}
	for i := range amount {
		if dirRight {
			d.ring = d.ring.Next()
		} else {
			d.ring = d.ring.Prev()
		}

		shouldIncrementPassword := d.ring.Value.(int) == 0 && ((d.passwordMethod == passwordMethodEnd && i+1 == amount) || (d.passwordMethod == passwordMethodClick))

		if shouldIncrementPassword {
			d.password += 1
		}
	}
	d.logger.Printf("The dial is rotated %s to point at %d\n", line, d.ring.Value.(int))
	return nil
}

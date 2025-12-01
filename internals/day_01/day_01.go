package day_01

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type Dial struct {
	logger   *log.Logger
	position int
	size     int
}

func NewDial() *Dial {
	return &Dial{
		logger:   log.New(os.Stdout, "", log.Lshortfile),
		position: 50,
		size:     100,
	}
}

func (d *Dial) GetPassword(reader io.Reader) (int, error) {
	d.logger.Printf("The dial starts at position %d\n", d.position)
	password := 0
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		err := d.turnDialUsingInstruction(line)
		if err != nil {
			return 0, err
		}
		if d.position == 0 {
			password += 1
		}
		d.logger.Printf("The dial is rotated %s to point at %d\n", line, d.position)
	}
	d.logger.Printf("The password is %d\n", password)
	return password, nil
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
	position := int(math.Remainder(float64(d.position+sign*amount), float64(d.size)))
	if position >= 0 {
		d.position = position
	} else {
		d.position = d.size + position
	}
	return nil
}

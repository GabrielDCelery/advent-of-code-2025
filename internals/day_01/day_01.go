package day_01

import (
	"bufio"
	"container/ring"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

const defaultRingSize = 100
const defaultRingPosition = 50

const (
	passwordMethodEnd   = 0x454E44
	passwordMethodClick = 0x434C49434B
)

type Dial struct {
	logger         *zap.Logger
	ring           *ring.Ring
	password       int
	passwordMethod int
}

func NewDial(passwordMethod string, logger *zap.Logger) (*Dial, error) {
	ringSize := defaultRingSize
	ringPosition := defaultRingPosition
	if logger == nil {
		logger = zap.NewNop()
	}
	logger = logger.With(zap.String("passwordMethod", passwordMethod))
	dial := &Dial{
		logger:         logger,
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
		return nil, fmt.Errorf("unhandled password method %s", passwordMethod)
	}
	return dial, nil
}

func (d *Dial) GetPassword(reader io.Reader) (int, error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		err := d.turnDialUsingInstruction(line)
		if err != nil {
			return 0, err
		}
	}
	d.logger.Info("password retrieved", zap.Int("password", d.password))
	return d.password, nil
}

func (d *Dial) turnDialUsingInstruction(line string) error {
	if len(line) == 0 {
		return fmt.Errorf("instruction can not be empty")
	}
	direction := line[:1]
	if direction != "L" && direction != "R" {
		return fmt.Errorf("unhandled instruction direction %s", line)
	}
	dirRight := direction == "R"
	amount, err := strconv.Atoi(line[1:])
	if err != nil {
		return fmt.Errorf("invalid amount in instruction %s", line)
	}
	for i := range amount {
		if dirRight {
			d.ring = d.ring.Next()
		} else {
			d.ring = d.ring.Prev()
		}

		isAtZero := d.ring.Value.(int) == 0
		isEndOfMove := (i + 1) == amount
		isEndMethod := isEndOfMove && d.passwordMethod == passwordMethodEnd
		isClickMethod := d.passwordMethod == passwordMethodClick
		shouldIncrementPassword := isAtZero && (isEndMethod || isClickMethod)

		if shouldIncrementPassword {
			d.password += 1
		}
	}
	d.logger.Debug(
		"dial rotated",
		zap.String("command", line),
		zap.Int("postion", d.ring.Value.(int)),
		zap.Int("password", d.password),
	)
	return nil
}

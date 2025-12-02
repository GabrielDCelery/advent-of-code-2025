package day02

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

const (
	ProductIDHasExactRepeat string = "exactrepeat"
	ProductIDHasAnyRepeat   string = "anyrepeat"
)

type Day2Solver struct {
	logger             *zap.Logger
	isProductIDInvalid isIDInvalidFunc
}

func NewDay2Solver(logger *zap.Logger, productValidator string) (*Day2Solver, error) {
	if logger == nil {
		logger = zap.NewNop()
	}
	day2Solver := &Day2Solver{
		logger:             logger,
		isProductIDInvalid: nil,
	}
	switch productValidator {
	case ProductIDHasExactRepeat:
		day2Solver.isProductIDInvalid = productIDHasExactRepeat
	case ProductIDHasAnyRepeat:
		day2Solver.isProductIDInvalid = productIDHasAnyRepeat
	default:
		return nil, fmt.Errorf("unhandled product validator %s", productValidator)
	}
	return day2Solver, nil
}

func (d *Day2Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	scanner := createProductIdInputScanner(reader)

	invalidIDSum := 0

	for scanner.Scan() {
		productIDRange := scanner.Text()
		min, max, err := convertProductIDRangeToMinMax(productIDRange)
		if err != nil {
			return 0, err
		}
		for i := min; i <= max; i++ {
			select {
			case <-ctx.Done():
				return 0, nil
			default:
				id := strconv.Itoa(i)
				if d.isProductIDInvalid(id) {
					d.logger.Debug("found invalid product ID",
						zap.Int("productID", i),
						zap.String("productIDRange", productIDRange),
					)
					invalidIDSum += i
				}
			}
		}
	}

	return invalidIDSum, nil
}

func createProductIdInputScanner(reader io.Reader) *bufio.Scanner {
	scanner := bufio.NewScanner(reader)

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF == true && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, ','); i >= 0 {
			return i + 1, bytes.TrimSpace(data[:i]), nil
		}
		if atEOF == true {
			return len(data), bytes.TrimSpace(data), nil
		}
		return 0, nil, nil
	})

	return scanner
}

func convertProductIDRangeToMinMax(productIDRange string) (int, int, error) {
	ids := strings.Split(productIDRange, "-")
	if len(ids) != 2 {
		return 0, 0, fmt.Errorf("failed to split product ID range %s", productIDRange)
	}
	min, err := strconv.Atoi(ids[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get minimum product ID from %s, reason: %v", productIDRange, err)
	}
	max, err := strconv.Atoi(ids[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get maximum product ID from %s, reason: %v", productIDRange, err)
	}
	return min, max, nil
}

type isIDInvalidFunc func(id string) bool

func productIDHasExactRepeat(id string) bool {
	sequenceLen := ((len(id) + 1) / 2)
	return isSequenceRepeating(id, sequenceLen)
}

func productIDHasAnyRepeat(id string) bool {
	for sequenceLen := 1; sequenceLen <= (len(id) / 2); sequenceLen++ {
		if isSequenceRepeating(id, sequenceLen) {
			return true
		}
	}
	return false
}

func isSequenceRepeating(data string, sequenceLen int) bool {
	if sequenceLen > (len(data) / 2) {
		return false
	}
	if (len(data) % sequenceLen) != 0 {
		return false
	}
	for i := 1; i < (len(data) / sequenceLen); i++ {
		for j := range sequenceLen {
			if data[j] != data[j+i*sequenceLen] {
				return false
			}
		}
	}
	return true
}

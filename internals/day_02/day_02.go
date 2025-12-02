package day_02

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type ProductIDInvalidSumOpt string

const (
	SomeSequenceRepeatedTwice        ProductIDInvalidSumOpt = "SOME_SEQUENCE_REPEATED_TWICE"
	SomeSequenceRepeatedAtleastTwice ProductIDInvalidSumOpt = "SOME_SEQUENCE_REPEATED_ATLEAST_TWICE"
)

type Day2Solver struct {
	logger                     *zap.Logger
	getInvalidProductIdSumFunc getInvalidProductIdSumFunc
}

func NewDay2Solver(logger *zap.Logger, productIdInvalidSumOpt ProductIDInvalidSumOpt) (*Day2Solver, error) {
	if logger == nil {
		logger = zap.NewExample()
	}
	day2Solver := &Day2Solver{
		logger:                     logger,
		getInvalidProductIdSumFunc: nil,
	}
	switch productIdInvalidSumOpt {
	case SomeSequenceRepeatedTwice:
		day2Solver.getInvalidProductIdSumFunc = getInvalidProductIdSumWhenSequenceIsRepeatedTwice
	case SomeSequenceRepeatedAtleastTwice:
		return nil, fmt.Errorf("unhandled option %s", productIdInvalidSumOpt)
	default:
		return nil, fmt.Errorf("unhandled option %s", productIdInvalidSumOpt)
	}
	return day2Solver, nil
}

func (d *Day2Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := createProductIdInputScanner(reader)

	invalidIDSum := 0

	for scanner.Scan() {
		productIDRange := scanner.Text()
		sum, err := d.getInvalidProductIdSumFunc(ctx, productIDRange)
		if err != nil {
			return 0, err
		}
		invalidIDSum += sum
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

type getInvalidProductIdSumFunc func(ctx context.Context, productIDRange string) (int, error)

func getInvalidProductIdSumWhenSequenceIsRepeatedTwice(ctx context.Context, productIDRange string) (int, error) {
	min, max, err := convertProductIDRangeToMinMax(productIDRange)
	if err != nil {
		return 0, err
	}
	sum := 0
	for i := min; i <= max; i++ {
		select {
		case <-ctx.Done():
		default:
			id := strconv.Itoa(i)
			if math.Remainder(float64(len(id)), 2) != 0 {
				continue
			}
			if id[:(len(id)/2)] == id[(len(id)/2):] {
				sum += i
			}
		}
	}
	return sum, nil
}

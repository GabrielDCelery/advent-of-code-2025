package day_02

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"math"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Day2Solver struct {
	logger *zap.Logger
}

func NewDay2Solver(logger *zap.Logger) *Day2Solver {
	if logger == nil {
		logger = zap.NewExample()
	}
	day2Solver := &Day2Solver{
		logger: logger,
	}
	return day2Solver
}

func (d *Day2Solver) Solve(ctx context.Context, reader io.Reader) (int, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scanner := createProductIdInputScanner(reader)

	invalidIDSum := 0

	for scanner.Scan() {
		productIDRange := scanner.Text()
		min, max, err := convertProductIDRangeToMinMax(productIDRange)
		if err != nil {
			return 0, err
		}
		invalidIDsChan := getInvalidIDs(ctx, min, max)
		for invalidID := range invalidIDsChan {
			invalidIDSum += invalidID
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

func convertProductIDRangeToMinMax(productIdRange string) (int, int, error) {
	ids := strings.Split(productIdRange, "-")
	min, err := strconv.Atoi(ids[0])
	if err != nil {
		return 0, 0, err
	}
	max, err := strconv.Atoi(ids[1])
	if err != nil {
		return 0, 0, err
	}
	return min, max, nil
}

func getInvalidIDs(ctx context.Context, min int, max int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := min; i <= max; i++ {
			select {
			case <-ctx.Done():
			default:
				id := strconv.Itoa(i)
				if math.Remainder(float64(len(id)), 2) != 0 {
					continue
				}
				if id[:(len(id)/2)] == id[(len(id)/2):] {
					out <- i
				}
			}
		}
	}()
	return out
}

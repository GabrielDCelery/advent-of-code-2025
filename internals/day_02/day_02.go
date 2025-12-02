package day_02

import (
	"bufio"
	"bytes"
	"io"

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

func (d *Day2Solver) Solve(reader io.Reader) (int, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF == true && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, ','); i >= 0 {
			return i + 1, data[:i], nil
		}
		if atEOF == true {
			return len(data), data, nil
		}
		return 0, nil, nil
	})
	for scanner.Scan() {
		item := scanner.Text()
		d.logger.Debug("reading line", zap.String("item", item))
	}
	return 0, nil
}

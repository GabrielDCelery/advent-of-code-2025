package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/GabrielDCelery/advent-of-code-2025/internals/day_02"
	"github.com/GabrielDCelery/advent-of-code-2025/internals/logging"
	"go.uber.org/zap"
)

func main() {
	validator, ok := os.LookupEnv("VALIDATOR")

	if !ok {
		validator = day_02.ProductIDHasExactRepeat
	}

	logLevel, ok := os.LookupEnv("LOGLEVEL")

	if !ok {
		logLevel = "info"
	}

	filePath := flag.String("file", "", "input file path")

	passwordMethod := flag.String("method", "end", "password method")

	flag.Parse()

	if *filePath == "" {
		log.Fatalf("missing param 'file', received %s", *filePath)
	}

	if *passwordMethod == "" {
		log.Fatalf("missing param 'method', received %s", *passwordMethod)
	}

	file, err := os.Open(*filePath)

	if err != nil {
		log.Fatalf("failed to open file at path %s", *filePath)
	}

	defer file.Close()

	logger := logging.NewLogger(logLevel)
	defer logger.Sync()

	day2Solver, err := day_02.NewDay2Solver(logger, validator)

	if err != nil {
		logger.Fatal("failed to instantiate day 2 problem solver", zap.Error(err))
	}

	solution, err := day2Solver.Solve(context.Background(), file)

	if err != nil {
		logger.Fatal("failed to run day 2 problem solver", zap.Error(err))
	}

	logger.Info("solved day 2 problem", zap.Int("solution", solution))
}

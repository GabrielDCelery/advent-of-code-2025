package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/GabrielDCelery/advent-of-code-2025/internals/day03"
	"github.com/GabrielDCelery/advent-of-code-2025/internals/logging"
	"go.uber.org/zap"
)

func main() {
	batteryCountStr := flag.String("batteryCount", "2", "the number of batteries in a power bank we want to calculate the joltage for")
	filePath := flag.String("file", "", "path to the input file containing product ID ranges")
	logLevel := flag.String("logLevel", "info", "log level for application")

	flag.Parse()

	if *filePath == "" {
		log.Fatalf("missing required flag: -file")
	}

	batteryCount, err := strconv.Atoi(*batteryCountStr)

	if err != nil {
		log.Fatalf("battery count has to be valid integer, but got %s", *batteryCountStr)
	}

	file, err := os.Open(*filePath)

	if err != nil {
		log.Fatalf("failed to open file at path %s: %v", *filePath, err)
	}

	defer file.Close()

	logger := logging.NewLogger(*logLevel)
	defer logger.Sync()

	day3Solver, err := day03.NewDay3Solver(batteryCount, logger)

	if err != nil {
		logger.Fatal("failed to instantiate day 3 problem solver", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	solution, err := day3Solver.Solve(ctx, file)

	if err != nil {
		logger.Fatal("failed to run day 3 problem solver", zap.Error(err))
	}

	logger.Info("solved day 3 problem", zap.Int("solution", solution))
}

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/GabrielDCelery/advent-of-code-2025/internals/day07"
	"github.com/GabrielDCelery/advent-of-code-2025/internals/logging"
	"go.uber.org/zap"
)

func main() {
	filePath := flag.String("file", "", "path to the input file containing product ID ranges")
	logLevel := flag.String("logLevel", "info", "log level for application")

	flag.Parse()

	if *filePath == "" {
		log.Fatalf("missing required flag: -file")
	}

	file, err := os.Open(*filePath)

	if err != nil {
		log.Fatalf("failed to open file at path %s: %v", *filePath, err)
	}

	defer file.Close()

	logger := logging.NewLogger(*logLevel)
	defer logger.Sync()

	solver, err := day07.NewDay7Solver(logger)

	if err != nil {
		logger.Fatal("failed to instantiate day 7 problem solver", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	solution, err := solver.Solve(ctx, file)

	if err != nil {
		logger.Fatal("failed to run day 7 problem solver", zap.Error(err))
	}

	logger.Info("solved day 7 problem", zap.Int("solution", solution))
}

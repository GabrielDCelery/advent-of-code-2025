package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/GabrielDCelery/advent-of-code-2025/internals/day04"
	"github.com/GabrielDCelery/advent-of-code-2025/internals/logging"
	"go.uber.org/zap"
)

func main() {
	removeRolls := flag.String("removeRolls", "n", "whether to remove rolls or not after we reach them (y or n)")
	filePath := flag.String("file", "", "path to the input file containing product ID ranges")
	logLevel := flag.String("logLevel", "info", "log level for application")

	flag.Parse()

	var shouldRemoveRolls bool
	switch *removeRolls {
	case "y":
		shouldRemoveRolls = day04.RemoveRolls
	case "n":
		shouldRemoveRolls = day04.DontRemoveRolls
	default:
		log.Fatalf("incorrect flag of '%s' for removeRolls, valid values are 'y' or 'n'", *removeRolls)
	}

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

	day3Solver, err := day04.NewDay4Solver(logger)

	if err != nil {
		logger.Fatal("failed to instantiate day 4 problem solver", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	solution, err := day3Solver.Solve(ctx, file, shouldRemoveRolls)

	if err != nil {
		logger.Fatal("failed to run day 4 problem solver", zap.Error(err))
	}

	logger.Info("solved day 4 problem", zap.Int("solution", solution))
}

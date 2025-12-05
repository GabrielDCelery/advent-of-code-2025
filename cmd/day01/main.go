package main

import (
	"flag"
	"log"
	"os"

	"github.com/GabrielDCelery/advent-of-code-2025/internals/day01"
	"github.com/GabrielDCelery/advent-of-code-2025/internals/logging"
	"go.uber.org/zap"
)

func main() {
	filePath := flag.String("file", "", "input file path")
	passwordMethod := flag.String("passwordMethod", "end", "password method (end or click)")
	logLevel := flag.String("logLevel", "info", "log level for application")

	flag.Parse()

	if *filePath == "" {
		log.Fatalf("missing flag: -file")
	}

	file, err := os.Open(*filePath)

	if err != nil {
		log.Fatalf("failed to open file at path %s", *filePath)
	}

	defer file.Close()

	logger := logging.NewLogger(*logLevel)
	defer logger.Sync()

	solver, err := day01.NewDay1Solver(*passwordMethod, logger)

	if err != nil {
		log.Fatalf("failed to instantiate day 1 puzzle solver %v", err)
	}

	solution, err := solver.Solve(file)

	if err != nil {
		log.Fatalf("failed to get password: %v", err)
	}

	logger.Info("solved day 1 puzzle", zap.Int("password", solution))
}

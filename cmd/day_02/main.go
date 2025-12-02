package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/GabrielDCelery/advent-of-code-2025/internals/day_02"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
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
	var atom zap.AtomicLevel
	switch logLevel {
	case "debug":
		atom = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	default:
		atom = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	encoderCfg := zap.NewProductionEncoderConfig()
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()
	day2Solver := day_02.NewDay2Solver(logger)
	solution, err := day2Solver.Solve(context.Background(), file)
	if err != nil {
		logger.Fatal("failed to run day 2 problem solver", zap.Error(err))
	}
	logger.Info("solved day 2 problem", zap.Int("solution", solution))
}

package main

import (
	"flag"
	"log"
	"os"

	"github.com/GabrielDCelery/advent-of-code-2025/internals/day_01"
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
	dial, err := day_01.NewDial(*passwordMethod, logger)
	if err != nil {
		log.Fatalf("failed to instantiate dial %v", err)
	}
	_, err = dial.GetPassword(file)
	if err != nil {
		log.Fatalf("failed to get password: %v", err)
	}
}

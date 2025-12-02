package main

import (
	"flag"
	"log"
	"os"

	"github.com/GabrielDCelery/advent-of-code-2025/internals/day_01"
	"github.com/GabrielDCelery/advent-of-code-2025/internals/logging"
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

	dial, err := day_01.NewDial(*passwordMethod, logger)

	if err != nil {
		log.Fatalf("failed to instantiate dial %v", err)
	}

	_, err = dial.GetPassword(file)

	if err != nil {
		log.Fatalf("failed to get password: %v", err)
	}
}

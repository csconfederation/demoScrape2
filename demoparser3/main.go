package main

import (
	"os"

	"github.com/csconfederation/demoparser3/logger"
)

func main() {
	demo := "path/to/demo"

	file, err := os.Open(demo)
	if err != nil {
		logger.Error("Failed to open demo file: %v", err)
	}
	defer file.Close()

	g, err := ProcessDemo(file)
	if err != nil {
		logger.Error("Failed to parse demo", "error", err)
	}

	logger.Info(g.CurrentRound.Planter.Name)
}

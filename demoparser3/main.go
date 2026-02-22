package main

import (
	"os"

	"github.com/csconfederation/demoparser3/logger"
)

func main() {
	demo := "demos/s18-M06-Magnolias-vs-TheMolotovs-mid7088-0_de_overpass-2025-10-24_02-09-12.dem"

	file, err := os.Open(demo)
	if err != nil {
		logger.Error("Failed to open demo file: %v", err)
	}
	defer file.Close()

	_, err = ProcessDemo(file)
	if err != nil {
		logger.Error("Failed to parse demo", "error", err)
	}
	logger.Info("Done processing demo")
}

package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // ensure different results on every run

	// Use DefaultResourceConfig() for auto-detection, or set manually:
	cfg := DefaultResourceConfig()

	NewSimulator(cfg).Run()
}

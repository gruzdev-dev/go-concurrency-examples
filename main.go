package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"race-lab/config"
	"race-lab/runner"
)

func main() {
	configPath := flag.String("config", "configs/race.json", "path to config")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	expected := int64(cfg.GoroutinesCount * cfg.Iterations)

	fmt.Printf("Mode: %s\n", cfg.Mode)
	fmt.Printf("Goroutines: %d\n", cfg.GoroutinesCount)
	fmt.Printf("Iterations: %d\n", cfg.Iterations)
	fmt.Printf("Expected: %d\n", expected)

	r := runner.New(cfg.GoroutinesCount, cfg.Iterations, cfg.Mode)

	start := time.Now()
	actual := r.Run()
	elapsed := time.Since(start)

	fmt.Printf("Actual: %d\n", actual)
	fmt.Printf("Duration: %v\n", elapsed)

	if actual != expected {
		fmt.Printf("WARNING: lost %d updates due to data race\n", expected-actual)
	}
}

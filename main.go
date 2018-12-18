package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yurishkuro/microsim/config"
	"github.com/yurishkuro/microsim/model"
)

var simulation = flag.String("c", "hotrod", "name of the simulation config or path to a JSON config file")
var printConfig = flag.Bool("o", false, "if present, print the config and exit")
var printValidated = flag.Bool("O", false, "if present, print the config with defaults and exit")
var duration = flag.Int("d", 10, "simulation duration in seconds")
var workers = flag.Int("w", 3, "number of workers (tests) to run in parallel")
var repeats = flag.Int("r", 0, "number of requests (repeats) each worker will send (default - as long as simulation is running)")

func main() {
	flag.Parse()

	if *simulation == "" {
		fmt.Fprintln(os.Stderr, "ERROR: simulation configuration name is required")
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.Get(*simulation)
	if err != nil {
		log.Fatalf("cannot load config %s: %v", *simulation, err)
	}

	// for now always print the config
	enc := json.NewEncoder(os.Stdout)
	// enc.SetIndent("", "  ")
	enc.Encode(cfg)
	if *printConfig {
		os.Exit(0)
	}

	r := &model.Registry{}
	r.RegisterServices(cfg.Services)

	if err := cfg.Validate(r); err != nil {
		log.Fatalf("validation failed: %v", err)
	}

	if *printValidated {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(cfg)
		os.Exit(0)
	}

	if err := cfg.Start(); err != nil {
		log.Fatalf("start failed: %v", err)
	}

	log.Printf("services started")
	time.Sleep(3 * time.Second)

	cfg.TestDuration = time.Duration(*duration) * time.Second
	cfg.TestRunners = *workers
	cfg.Repeats = *repeats
	cfg.Run()

	cfg.Stop()
}

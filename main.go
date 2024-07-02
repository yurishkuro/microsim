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

var (
	simulation     = flag.String("c", "hotrod", "name of the simulation config or path to a JSON config file")
	printConfig    = flag.Bool("o", false, "if present, print the config and exit")
	printValidated = flag.Bool("O", false, "if present, print the config with defaults and exit")
	duration       = flag.Duration("d", 10*time.Second, "simulation duration")
	workers        = flag.Int("w", 3, "number of workers (tests) to run in parallel")
	repeats        = flag.Int("r", 0, "number of requests (repeats) each worker will send (default 0, i.e. as long as simulation is running)")
	sleep          = flag.Duration("s", 100*time.Millisecond, "how long each worker sleeps between requests, as a way of controlling QPS")
)

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
	_ = enc.Encode(cfg)
	if *printConfig {
		os.Exit(0)
	}

	r := &model.Registry{}
	_ = r.RegisterServices(cfg.Services)

	if err := cfg.Validate(r); err != nil {
		log.Fatalf("validation failed: %v", err)
	}

	if *printValidated {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(cfg)
		os.Exit(0)
	}

	if err := cfg.Start(); err != nil {
		log.Fatalf("start failed: %v", err)
	}

	log.Printf("services started")
	time.Sleep(3 * time.Second)

	cfg.TestName = *simulation
	cfg.TestDuration = *duration
	cfg.TestRunners = *workers
	cfg.Repeats = *repeats
	cfg.SleepBetweenRequests = *sleep
	cfg.Run()

	cfg.Stop()
}

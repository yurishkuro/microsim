package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/yurishkuro/microsim/model"
)

var duration = flag.Int("s", 10, "simulation duration in seconds")
var workers = flag.Int("w", 3, "number of workers (tests) to run in parallel")
var repeats = flag.Int("r", 0, "number of requests (repeats) each worker will send (default - as long as simulation is running)")

func main() {
	flag.Parse()

	cfg := hotrod
	out, err := json.Marshal(&cfg)
	if err != nil {
		panic(err.Error())
	}

	os.Stdout.WriteString(string(out) + "\n")

	r := &model.Registry{}
	r.RegisterServices(cfg.Services)

	if err := cfg.Validate(r); err != nil {
		log.Fatalf("validation failed: %v", err)
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

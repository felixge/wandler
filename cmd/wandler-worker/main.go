package main

import (
	"flag"
	"fmt"
	"github.com/felixge/wandler/config"
	"github.com/felixge/wandler/queue"
	"github.com/felixge/wandler/log"
	"github.com/felixge/wandler/worker"
	"os"
)

var DefaultConfig = Config{
	LogLevel:      "debug",
	LogTimeFormat: "15:04:05.999",
	JobQueue:      "redis://localhost/",
}

type Config struct {
	LogLevel      string
	LogTimeFormat string
	JobQueue      string
}

func main() {
	var (
		flags      = flag.NewFlagSet("wandler-worker", flag.ExitOnError)
		configPath = flags.String("f", "", "Config file to load.")
	)

	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conf := DefaultConfig
	if *configPath != "" {
		if err := config.ReadFile(*configPath, &conf); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	log, err := log.NewLogger(conf.LogLevel, conf.LogTimeFormat, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Notice("Starting wandler worker")

	log.Debug("Creating job queue: %s", conf.JobQueue)
	q, err := queue.NewQueue(conf.JobQueue, log)
	if err != nil {
		log.Emergency("Could not create job queue: %s", err)
	}

	w, err := worker.NewWorker(log, q)
	if err := w.Run(); err != nil {
		os.Exit(1)
	}
}

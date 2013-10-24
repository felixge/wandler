package main

import (
	"flag"
	"fmt"
	"github.com/felixge/wandler/config"
	"github.com/felixge/wandler/http"
	"github.com/felixge/wandler/log"
	"github.com/felixge/wandler/queue"
	"net"
	gohttp "net/http"
	"os"
	"strings"
	"sync"
)

var DefaultConfig = Config{
	LogLevel:      "debug",
	LogTimeFormat: "15:04:05.999",
	HttpAddr:      ":8080",
	Queue:         "redis://localhost/",
	HttpURL:       "http://localhost:8080/",
}

type Config struct {
	LogLevel      string
	LogTimeFormat string
	HttpAddr      string
	Queue         string
	HttpURL       string
}

func main() {
	var (
		flags      = flag.NewFlagSet("wandler-server", flag.ExitOnError)
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

	// remove trailing /
	conf.HttpURL = strings.TrimSuffix(conf.HttpURL, "/")

	log, err := log.NewLogger(conf.LogLevel, conf.LogTimeFormat, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Notice("Starting wandler server")

	log.Debug("Creating job queue: %s", conf.Queue)
	q, err := queue.NewQueue(conf.Queue, log)
	if err != nil {
		log.Emergency("Could not create job queue: %s", err)
	}

	log.Debug("Creating http listener addr=%s", conf.HttpAddr)
	httpListener, err := net.Listen("tcp", conf.HttpAddr)
	if err != nil {
		log.Emergency("Could not create http listener: %s", err)
	}

	log.Debug("Creating http handler")
	httpHandler, err := http.NewHandler(http.HandlerConfig{
		Log:     log,
		Queue:   q,
		HttpURL: conf.HttpURL,
	})
	if err != nil {
		log.Emergency("Could not create http handler: %s", err)
	}

	log.Debug("Creating http server")
	httpServer := &gohttp.Server{Handler: httpHandler}

	var shutdown sync.WaitGroup

	shutdown.Add(1)
	go func() {
		defer shutdown.Done()
		log.Notice("Serving http clients now")
		if err := httpServer.Serve(httpListener); err != nil {
			log.Emergency("Error serving http: %s", err)
		}
	}()

	shutdown.Wait()
	log.Notice("Shutting down")
}

package main

import (
	"flag"
	"fmt"
	"github.com/felixge/wandler/config"
	"github.com/felixge/wandler/log"
	"net"
	"net/http"
	"os"
	"sync"
)

var DefaultConfig = Config{
	LogLevel:      "debug",
	LogTimeFormat: "15:04:05.999",
	HttpAddr:      ":8080",
}

type Config struct {
	LogLevel      string
	LogTimeFormat string
	HttpAddr      string
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

	log, err := log.NewLogger(conf.LogLevel, conf.LogTimeFormat, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.Notice("Starting wandler server")

	log.Debug("Creating http listener at: %s", conf.HttpAddr)
	httpListener, err := net.Listen("tcp", conf.HttpAddr)
	if err != nil {
		log.Emergency("Could not create http listener: %s", err)
	}

	log.Debug("Initializing http server")
	httpServer := &http.Server{}

	var shutdown sync.WaitGroup

	shutdown.Add(1)
	go func() {
		defer shutdown.Done()
		log.Debug("Serving http clients now")
		if err := httpServer.Serve(httpListener); err != nil {
			log.Emergency("Error serving http: %s", err)
		}
	}()

	log.Notice("Started wandler server")
	shutdown.Wait()
	log.Notice("Completed shutdown")
}

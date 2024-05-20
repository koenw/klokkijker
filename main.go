package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/alecthomas/kong"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
)

var CLI struct {
	Count      int    `long:"count" short:"c" help:"Amount of requests to send to each server every interfal." default:"3"`
	Interval   int    `long:"interval" short:"i" help:"Time in seconds to wait between NTP requests" default:"11"`
	Database   string `long:"database" short:"d" help:"Sqlite database to write results to for logkeeping."`
	Format     string `long:"format" short:"f" help:"Output format, one of 'json', 'cli', 'text'."`
	Prometheus bool   `long:"prometheus" help:"Export prometheus metrics"`
	ListenAddress string `long:"listen" help:"Address to export prometheus metrics on" default:":8123"`
	Servers []string `arg:"" optional:"" name:"servers" help:"NTP Servers." type:"string"`
}

func main() {
	kong.Parse(&CLI)

	switch CLI.Format {
	case "json":
		log.SetHandler(json.New(os.Stdout))
	case "text":
		log.SetHandler(text.New(os.Stdout))
	default:
		log.SetHandler(cli.New(os.Stderr))
	}

	if CLI.Prometheus {
		var wg sync.WaitGroup
		defer wg.Wait()
		wg.Add(1)
		handler := promHandler(CLI.Servers, CLI.Count, CLI.Interval)

		http.Handle("/metrics", handler)
		log.Info("Serving prometheus metrics on port 8123")
		http.ListenAndServe(CLI.ListenAddress, nil)
	} else {
		ch := make(chan NTPResponse)
		for _, server := range CLI.Servers {
			go ping(server, CLI.Count, ch)
		}
		for i := 0; i < (CLI.Count * len(CLI.Servers)); i++ {
			<-ch
		}
	}
}

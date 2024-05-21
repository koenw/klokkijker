package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/koenw/klokkijker/internal/ntp"
	"github.com/koenw/klokkijker/internal/prometheus"
	"github.com/spf13/cobra"
)

var (
	monitorCmd = &cobra.Command{
		Use:   "monitor <NTP Servers>",
		Short: "Continuously send NTP requests and export prometheus metrics",
		Long:  "Continuously send NTP requests to the given NTP servers and export prometheus\nmetrics over HTTP.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, servers []string) {
			setupLogging(outputFormat)

			intervalFloat, err := strconv.ParseFloat(interval, 64)
			if err != nil {
				log.Fatalf("Failed to parse %s into number", interval)
			}
			rpm := int(60 / intervalFloat)
			rampPeriod, _ := time.ParseDuration("0s")
			resps := ntp.GenerateRequests(servers, rpm, rampPeriod)

			handler := prometheus.Monitor(servers, resps)
			http.Handle("/metrics", handler)

			listenAddr := fmt.Sprintf("%s:%d", prometheusAddr, prometheusPort)

			log.Info(fmt.Sprintf("Serving metrics at http://%s/metrics", listenAddr))
			http.ListenAndServe(listenAddr, nil)
		},
	}

	count          int
	interval       string
	prometheusAddr string
	prometheusPort int
)

func init() {
	rootCmd.AddCommand(monitorCmd)

	monitorCmd.Flags().IntVarP(
		&count,
		"count", "c",
		1, "How many NTP requests to send each salvo")

	monitorCmd.Flags().StringVarP(
		&interval,
		"interval", "i",
		"1", "Interval between NTP request salvo's")

	monitorCmd.Flags().StringVarP(
		&prometheusAddr,
		"prometheus-address", "",
		"127.0.0.1", "Hostname or IP address for the prometheus exporter to listen on")

	monitorCmd.Flags().IntVarP(
		&prometheusPort,
		"prometheus-port", "",
		8123,
		"Port for the prometheus exporter to listen on")
}

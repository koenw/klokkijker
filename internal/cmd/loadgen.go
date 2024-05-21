package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/koenw/klokkijker/internal/ntp"
	"github.com/koenw/klokkijker/internal/prometheus"
	"github.com/spf13/cobra"
)

var (
	loadgenCmd = &cobra.Command{
		Use:   "loadgen <NTP Servers>",
		Short: "Generate NTP requests for load testing and export prometheus metrics (experimental)",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, servers []string) {
			setupLogging(outputFormat)
			ramp, err := time.ParseDuration(rampPeriod)
			if err != nil {
				log.Fatalf(fmt.Sprintf("Failed to parse %s into duration: %s", rampPeriod, err))
			}

			resps := ntp.GenerateRequests(servers, rpm, ramp)

			handler := prometheus.Monitor(servers, resps)
			http.Handle("/metrics", handler)

			listenAddr := fmt.Sprintf("%s:%d", prometheusAddr, prometheusPort)

			log.Info(fmt.Sprintf("Serving metrics at http://%s/metrics", listenAddr))
			http.ListenAndServe(listenAddr, nil)
		},
	}

	rpm                   int
	rampPeriod            string
	loadgenPrometheusAddr string
	loadgenPrometheusPort int
)

func init() {
	rootCmd.AddCommand(loadgenCmd)

	loadgenCmd.Flags().IntVarP(
		&rpm,
		"rpm", "",
		60, "Requests Per Minute to send to the NTP servers")

	loadgenCmd.Flags().StringVarP(
		&rampPeriod,
		"ramp-period", "",
		"0s", "Period during which to build-up our Requests-Per-Minute until we hit our target RPM")

	loadgenCmd.Flags().StringVarP(
		&loadgenPrometheusAddr,
		"prometheus-address", "",
		"127.0.0.1", "Hostname or IP address for the prometheus exporter to listen on")

	loadgenCmd.Flags().IntVarP(
		&loadgenPrometheusPort,
		"prometheus-port", "",
		8123,
		"Port for the prometheus exporter to listen on")
}

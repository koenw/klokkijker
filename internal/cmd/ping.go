package cmd

import (
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/koenw/klokkijker/internal/ntp"
	"github.com/spf13/cobra"
)

var (
	pingCmd = &cobra.Command{
		Use:   "ping [--count 1] <NTP Servers>",
		Short: "Send NTP requests and print the response",
		Long:  `Send NTP requests and print the response.`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, servers []string) {
			setupLogging(outputFormat)

			var resps chan ntp.NTPResponse

			if pingInterval != "" {
				intervalFloat, err := strconv.ParseFloat(pingInterval, 64)
				if err != nil {
					log.Fatalf("Failed to parse interval %s into number", pingInterval)
				}

				rpm := int(60 / intervalFloat)
				rampPeriod, _ := time.ParseDuration("0s")
				resps = ntp.GenerateRequests(servers, rpm, rampPeriod)

				for _ = range resps {
				}
			} else {
				resps = make(chan ntp.NTPResponse)

				for _, server := range servers {
					go ntp.Ping(server, pingCount, resps)
				}

				for i := 0; i < (pingCount * len(servers)); i++ {
					<-resps
				}
			}
		},
	}

	pingCount    int
	pingInterval string
)

func init() {
	rootCmd.AddCommand(pingCmd)

	pingCmd.Flags().IntVarP(
		&pingCount,
		"count", "c",
		1, "How many NTP requests to send")

	pingCmd.Flags().StringVarP(
		&pingInterval,
		"interval", "i",
		"", "Continuously send NTP requests, waiting <interval> seconds between requests")
}

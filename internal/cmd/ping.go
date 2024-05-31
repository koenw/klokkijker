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
		Example: `  klokkijker ntp.example.com

  # Same as above: send a single NTP request to ntp.example.com
  klokkijker ping ntp.example.com

  # Send 3 NTP requests to ntp.example.com
  klokkijker ping --count 3 ntp.example.com

  # Continuously send NTP requests to ntp2.example.com & ntp3.example.com,
  # once every 0.2 seconds.
  klokkijker ping --interval 0.2 ntp2.example.com ntp3.example.com

  # Continuously send NTP requests to ntp4.example.com, in batches of 5
  # every 1 second. Print the output in json format and pipe it through
  # *jq* to only print the Round Trip Time (in nanoseconds)
  klokkijker ping --format json --interval 1 --count 5 ntp4.example.com |jq .fields.rtt`,
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

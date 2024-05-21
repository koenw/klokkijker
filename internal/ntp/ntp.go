package ntp

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/beevik/ntp"

	"sync"
)

type NTPResponse struct {
	*ntp.Response
	Server string
	Err    error
}

func pingSingle(server string) (NTPResponse, error) {
	resp, err := ntp.Query(server)
	if err != nil {
		log.WithFields(log.Fields{
			"server": server,
		}).WithError(err)

	} else {
		log.WithFields(log.Fields{
			"time":           resp.Time,
			"stratum":        resp.Stratum,
			"referenceTime":  resp.ReferenceTime,
			"referenceID":    resp.ReferenceID,
			"offset":         resp.ClockOffset,
			"rtt":            resp.RTT,
			"precision":      resp.Precision,
			"rootDelay":      resp.RootDelay,
			"rootDispersion": resp.RootDispersion,
			"rootDistance":   resp.RootDistance,
			"poll":           resp.Poll,
			"server":         server,
		}).Info(server)
	}

	return NTPResponse{
		resp,
		server,
		err,
	}, err

}

func Ping(server string, count int, out chan NTPResponse) {
	var wg sync.WaitGroup

	for c := 0; c < count; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := pingSingle(server)
			if err != nil {
				log.WithFields(log.Fields{"server": server}).Errorf("%s", err)
			}
			out <- resp
		}()
	}

	wg.Wait()
}

// rpmToInterval takes an Requests Per Minutes (or Ops per Minute) and returns
// the interval between each occurence as a time.Duration
func rpmToInterval(rpm int) time.Duration {
	rpd := rpm * 60 * 24

	seconds := (float64(60*60*24) / float64(rpd))
	duration, err := time.ParseDuration(fmt.Sprintf("%fs", seconds))
	if err != nil {
		panic(fmt.Sprintf("Failed to parse duration from \"%fs\": %s", seconds, err))
	}
	return duration
}

// GenerateRequests generates NTP requests to the given servers at the given
// Requests Per Minute (RPM) building up from 0 rpm to the given rpm over the
// given ramp period in seconds.
func GenerateRequests(servers []string, rpm int, rampPeriod time.Duration) chan NTPResponse {
	resps := make(chan NTPResponse)

	for _, server := range servers {
		go func(server string, rpm int, rampPeriod time.Duration) {
			start := time.Now()
			rampDoneTime := start.Add(rampPeriod)
			sleepPeriod := rpmToInterval(rpm)

			var nextSleepPeriod time.Duration
			rampDone := false

			for {
				Ping(server, 1, resps)
				if !rampDone {
					now := time.Now()
					if now.Compare(rampDoneTime) >= 0 {
						rampDone = true
						continue
					}

					rampMillisLeft := rampDoneTime.UnixMilli() - now.UnixMilli()
					rampMillisTotal := rampDoneTime.UnixMilli() - start.UnixMilli()

					perc := 10 * float64(rampMillisLeft) / float64(rampMillisTotal)
					nextSleepPeriodMillis := perc * float64(sleepPeriod.Milliseconds())
					nextSleepPeriod, _ = time.ParseDuration(fmt.Sprintf("%fms", nextSleepPeriodMillis))
					if nextSleepPeriod.Milliseconds() > sleepPeriod.Milliseconds() {
						nextSleepPeriod = sleepPeriod
						rampDone = true
					}
					time.Sleep(nextSleepPeriod)
				} else {
					time.Sleep(sleepPeriod)
				}
			}
		}(server, rpm, rampPeriod)
	}

	return resps
}

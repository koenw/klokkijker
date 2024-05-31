package ntp

import (
	"fmt"
	"time"
	"math"

	"github.com/apex/log"
	"github.com/beevik/ntp"

	"sync"
)

type NTPResponse struct {
	*ntp.Response
	Server string
	Err    error
}


// PingSingle sends a single NTP Request to the given server.
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


// Ping sends `count` number of NTP requests to the given server, sending the
// response to the given `out` channel.
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
	// I'm increasing both to requests / day during the calculation to increase
	// resolution
	seconds := (float64(60*60*24) / float64(rpm * 60 * 24))
	var duration time.Duration
	var err error
	if seconds == math.Inf(1) {
		duration, err = time.ParseDuration(fmt.Sprintf("%ds", 1))
	} else {
		duration, err = time.ParseDuration(fmt.Sprintf("%fs", seconds))
	}
	if err != nil {
		panic(fmt.Sprintf("Failed to parse duration from \"%fs\": %s", seconds, err))
	}
	return duration
}


// GenerateRequests generates NTP requests to the given server at the given
// Requests Per Minute (RPM) building up from 0 rpm to the given rpm over the
// given ramp period in seconds, sending the replies to the supplied channel.
func generateRequests(server string, rpm int, rampPeriod time.Duration, resps chan NTPResponse) {
	rampSeconds := rampPeriod.Seconds()
	requestInterval := rpmToInterval(rpm)

	for i := 0; float64(i) < rampSeconds; i++ {
		nextRPM := (float64(i) / rampSeconds) * float64(rpm)
		nextRequestInterval := rpmToInterval(int(nextRPM))
		log.Debug(fmt.Sprintf("Setting RPM to %f (interval %s)", nextRPM, nextRequestInterval))
		count := 1_000_000_000 / nextRequestInterval.Nanoseconds()
		for j := int64(0); j < count; j++ {
			Ping(server, 1, resps)
			time.Sleep(nextRequestInterval)
		}
	}

	log.Debug(fmt.Sprintf("Setting RPM to %d (interval %s)", rpm, requestInterval))
	for {
		Ping(server, 1, resps)
		time.Sleep(requestInterval)
	}
}

// GenerateRequests generates NTP requests to the given servers at the given
// Requests Per Minute (RPM) building up from 0 rpm to the given rpm over the
// given ramp period in seconds and returns a channel of `NTPResponse`s.
func GenerateRequests(servers []string, rpm int, rampPeriod time.Duration) chan NTPResponse {
	resps := make(chan NTPResponse)

	for _, server := range servers {
		go generateRequests(server, rpm, rampPeriod, resps)
	}

	return resps
}

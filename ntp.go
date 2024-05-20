package main

import (
	"github.com/apex/log"
	"github.com/beevik/ntp"

	"sync"
)

type NTPResponse struct {
	*ntp.Response
	server string
	err    error
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

func ping(server string, count int, out chan<- NTPResponse) {
	var wg sync.WaitGroup

	for c := 0; c < count; c++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := pingSingle(server)
			if err != nil {
				log.WithFields(log.Fields{"server": server}).WithError(err)
			}
			out <- resp
		}()
	}

	wg.Wait()
}

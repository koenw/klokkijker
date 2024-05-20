package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func collectMetrics(server string, metrics *metrics, count int, interval int) {
	resps := make(chan NTPResponse)

	for {
		for _, server := range CLI.Servers {
			go ping(server, CLI.Count, resps)
		}

		for i := 0; i < (CLI.Count * len(CLI.Servers)); i++ {
			resp := <-resps

			metrics.requestsTotal.Inc()

			if resp.err != nil {
				metrics.requestsError.Inc()
			} else {
				metrics.offset.Set(float64(resp.ClockOffset.Nanoseconds()))
				metrics.rtt.Set(float64(resp.RTT.Nanoseconds()))
				metrics.precision.Set(float64(resp.Precision.Nanoseconds()))
				metrics.rootDelay.Set(float64(resp.RootDelay.Nanoseconds()))
				metrics.rootDispersion.Set(float64(resp.RootDispersion.Nanoseconds()))
				metrics.rootDistance.Set(float64(resp.RootDistance.Nanoseconds()))
				metrics.minError.Set(float64(resp.MinError.Nanoseconds()))
				metrics.leap.Set(float64(resp.Leap))
				metrics.poll.Set(float64(resp.Poll))
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

type metrics struct {
	offset         prometheus.Gauge
	rtt            prometheus.Gauge
	precision      prometheus.Gauge
	rootDelay      prometheus.Gauge
	rootDispersion prometheus.Gauge
	rootDistance   prometheus.Gauge
	requestsTotal  prometheus.Counter
	requestsError  prometheus.Counter
	leap           prometheus.Gauge
	minError       prometheus.Gauge
	poll           prometheus.Gauge
}

func newMetrics(reg prometheus.Registerer, server string) *metrics {
	namespace := "klokkijker"
	m := &metrics{
		offset: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "offset_nanoseconds",
			Help:        "Estimated offset of the local clock relative to the server's clock",
		}),
		rtt: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "rtt_nanoseconds",
			Help:        "Measured Round-Trip-Time between the client and the server",
		}),
		precision: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "precision_nanoseconds",
			Help:        "Reported precision of the server's clock",
		}),
		rootDelay: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "rootDelay_nanoseconds",
			Help:        "Server's estimated aggregate round-trip-time to the stratum 1 server",
		}),
		rootDispersion: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "rootDispersion_nanoseconds",
			Help:        "Server's estimated maximum measurement error relative to the stratum 1 server",
		}),
		rootDistance: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "rootDistance_nanoseconds",
			Help:        "Estimate of the total synchronization distance between the client and the stratum 1 server",
		}),
		requestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "requests_total",
			Help:        "Total amount of requests send to this server",
		}),
		requestsError: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "requests_error",
			Help:        "Amount of requests without valid response",
		}),
		leap: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "leap",
			Help:        "Indication if a leap second should be inserted or deleted in the last minute of the current month",
		}),
		minError: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "minError_nanoseconds",
			Help:        "Estimated lower bound on the error between server and client clocks",
		}),
		poll: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   namespace,
			ConstLabels: prometheus.Labels{"server": server},
			Name:        "poll",
			Help:        "Indication of when to next send a NTP request",
		}),
	}

	err := reg.Register(m.offset)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.rtt)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.precision)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.rootDelay)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.rootDispersion)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.rootDistance)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.requestsTotal)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.requestsError)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.leap)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.minError)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	err = reg.Register(m.poll)
	if err != nil {
		log.Error(fmt.Sprintf("%s", err))
		log.WithError(err)
	}

	return m
}

func promHandler(servers []string, count int, interval int) http.Handler {
	reg := prometheus.NewRegistry()

	for _, server := range servers {
		m := newMetrics(reg, server)
		go collectMetrics(server, m, count, interval)
	}

	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}

func recordMetrics(chan<- NTPResponse) {
}

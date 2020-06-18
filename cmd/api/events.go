package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/signalcd/signalcd/signalcd"
)

// Events to Deployments that should be sent via SSE (Server Sent Events)
type Events interface {
	SubscribeDeployments(channel chan signalcd.Deployment) signalcd.Subscription
	UnsubscribeDeployments(s signalcd.Subscription)
}

func deploymentEventsHandler(logger log.Logger, r *prometheus.Registry, events Events) func(w http.ResponseWriter, r *http.Request) {
	labels := prometheus.Labels{"events": "deployments"}
	eventDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "http_event_duration_seconds",
		Help:        "Duration and error code for server sent events",
		ConstLabels: labels,
	}, []string{"status"})

	subscribers := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "http_event_subscriptions",
		Help:        "Number of current subscribers",
		ConstLabels: labels,
	})

	r.MustRegister(eventDuration, subscribers)

	observeEvent := func(duration time.Duration, err error) {
		if err != nil {
			level.Warn(logger).Log(
				"msg", "failed to send server sent event",
				"err", err,
				logLabels(labels),
			)
			eventDuration.WithLabelValues("error").Observe(duration.Seconds())
		} else {
			level.Debug(logger).Log(
				"msg", "successfully sent server sent event",
				logLabels(labels),
			)
			eventDuration.WithLabelValues("success").Observe(duration.Seconds())
		}
	}

	observeSubscription := func() {
		subscribers.Inc()
		level.Debug(logger).Log("msg", "subscriber to events")
	}
	observeUnsubscription := func() {
		subscribers.Dec()
		level.Debug(logger).Log("msg", "unsubscriber from events")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, `{"message":"server sent events unsupported"}`, http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		deploymentEvents := make(chan signalcd.Deployment, 8)
		subscriptions := events.SubscribeDeployments(deploymentEvents)
		observeSubscription()

		defer func() {
			events.UnsubscribeDeployments(subscriptions)
			observeUnsubscription()
		}()

		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				close(deploymentEvents)
				return
			case deployment := <-deploymentEvents:
				start := time.Now()

				payload, err := json.Marshal(deploymentOpenAPI(deployment))
				if err != nil {
					observeEvent(time.Since(start), err)
					continue
				}

				_, err = fmt.Fprintf(w, "data: %s\n\n", payload)
				if err != nil {
					observeEvent(time.Since(start), err)
					continue
				}

				flusher.Flush()

				observeEvent(time.Since(start), nil)
			}
		}
	}
}

func logLabels(labels prometheus.Labels) []string {
	var s []string
	for k, v := range labels {
		s = append(s, k)
		s = append(s, v)
	}
	return s
}

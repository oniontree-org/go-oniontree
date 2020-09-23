package evtmetrics_test

import (
	"context"
	"github.com/oniontree-org/go-oniontree/scanner"
	"github.com/oniontree-org/go-oniontree/scanner/evtmetrics"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestMetrics_ReadEvents(t *testing.T) {
	eventCh := make(chan scanner.Event)
	eventCh2 := make(chan scanner.Event)
	emitEvent := func(event scanner.Event) {
		select {
		case eventCh <- event:
		}
		// Use sleep so that the cache has enough time to change its internal state.
		time.Sleep(400 * time.Millisecond)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metrics := &evtmetrics.Metrics{}
	metrics2 := &evtmetrics.Metrics{}
	exitCh := make(chan struct{}, 2)

	go func() {
		if err := metrics.ReadEvents(ctx, eventCh, eventCh2); err != nil {
			log.Printf("%s\n", err)
			return
		}
		exitCh <- struct{}{}
	}()
	go func() {
		if err := metrics2.ReadEvents(ctx, eventCh2, nil); err != nil {
			log.Printf("%s\n", err)
			return
		}
		exitCh <- struct{}{}
	}()

	serviceID := "oniontree"
	addresses := []string{
		"http://onions52ehmf4q75.onion",
		"http://onions53ehmf4q75.onion",
	}
	directory := "./test-dir"

	//
	// Bring the first address online
	//
	func(metrics ...*evtmetrics.Metrics) {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOnline,
			URL:       addresses[0],
			ServiceID: serviceID,
			Directory: directory,
		})

		for _, metrics := range metrics {
			gauge, err := metrics.Get().GetMetricWithLabelValues(serviceID, addresses[0], directory)

			if !assert.NoError(t, err) {
				t.Fatal(err)
			}

			metric := io_prometheus_client.Metric{}
			err = gauge.Write(&metric)

			if !assert.NoError(t, err) {
				t.Fatal(err)
			}

			if !assert.Equal(t, float64(scanner.StatusOnline), metric.GetGauge().GetValue()) {
				t.Fatal("unexpected value")
			}
		}
	}(metrics, metrics2)

	//
	// Bring the first address offline
	//
	func(metrics ...*evtmetrics.Metrics) {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOffline,
			URL:       addresses[0],
			ServiceID: serviceID,
			Directory: directory,
		})

		for _, metrics := range metrics {
			gauge, err := metrics.Get().GetMetricWithLabelValues(serviceID, addresses[0], directory)

			if !assert.NoError(t, err) {
				t.Fatal(err)
			}

			metric := io_prometheus_client.Metric{}
			err = gauge.Write(&metric)

			if !assert.NoError(t, err) {
				t.Fatal(err)
			}

			if !assert.Equal(t, float64(scanner.StatusOffline), metric.GetGauge().GetValue()) {
				t.Fatal("unexpected value")
			}
		}
	}(metrics, metrics2)

	close(eventCh)

	for i := 0; i < cap(exitCh); i++ {
		<-exitCh
	}
}

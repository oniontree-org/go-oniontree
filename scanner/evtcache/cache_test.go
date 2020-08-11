package evtcache_test

import (
	"context"
	"github.com/onionltd/go-oniontree/scanner"
	"github.com/onionltd/go-oniontree/scanner/evtcache"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestCache_ReadEvents(t *testing.T) {
	eventCh := make(chan scanner.Event)
	emitEvent := func(event scanner.Event) {
		select {
		case eventCh <- event:
		}
		// Use sleep so that the cache has enough time to change its internal state.
		time.Sleep(400 * time.Millisecond)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cache := &evtcache.Cache{}
	exitCh := make(chan struct{})

	go func() {
		if err := cache.ReadEvents(ctx, eventCh); err != nil {
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

	//
	// Bring the first address online
	//
	func() {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOnline,
			URL:       addresses[0],
			ServiceID: serviceID,
		})
		addrs, ok := cache.GetAddresses(serviceID)

		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, map[string]scanner.Status{
			addresses[0]: scanner.StatusOnline,
		}, addrs) {
			t.Fatal("invalid addresses returned")
		}

		onlineAddrs, ok := cache.GetOnlineAddresses(serviceID)
		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, []string{addresses[0]}, onlineAddrs) {
			t.Fatal("invalid online addresses returned")
		}
	}()

	//
	// Bring the second address online
	//
	func() {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOnline,
			URL:       addresses[1],
			ServiceID: serviceID,
		})
		addrs, ok := cache.GetAddresses(serviceID)

		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, map[string]scanner.Status{
			addresses[0]: scanner.StatusOnline,
			addresses[1]: scanner.StatusOnline,
		}, addrs) {
			t.Fatal("invalid addresses returned")
		}

		onlineAddrs, ok := cache.GetOnlineAddresses(serviceID)
		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, []string{addresses[0], addresses[1]}, onlineAddrs) {
			t.Fatal("invalid online addresses returned")
		}
	}()

	//
	// Bring the first address offline
	//
	func() {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOffline,
			URL:       addresses[0],
			ServiceID: serviceID,
		})
		addrs, ok := cache.GetAddresses(serviceID)

		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, map[string]scanner.Status{
			addresses[0]: scanner.StatusOffline,
			addresses[1]: scanner.StatusOnline,
		}, addrs) {
			t.Fatal("invalid addresses returned")
		}

		onlineAddrs, ok := cache.GetOnlineAddresses(serviceID)
		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, []string{addresses[1]}, onlineAddrs) {
			t.Fatal("invalid online addresses returned")
		}
	}()

	//
	// Bring the second address offline
	//
	func() {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOffline,
			URL:       addresses[1],
			ServiceID: serviceID,
		})
		addrs, ok := cache.GetAddresses(serviceID)

		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, map[string]scanner.Status{
			addresses[0]: scanner.StatusOffline,
			addresses[1]: scanner.StatusOffline,
		}, addrs) {
			t.Fatal("invalid addresses returned")
		}

		onlineAddrs, ok := cache.GetOnlineAddresses(serviceID)
		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, []string{}, onlineAddrs) {
			t.Fatal("invalid online addresses returned")
		}
	}()

	//
	// Send WorkerStopped event, this should remove the first address from the map.
	//
	func() {
		emitEvent(scanner.WorkerStopped{
			URL:       addresses[0],
			ServiceID: serviceID,
		})
		addrs, ok := cache.GetAddresses(serviceID)

		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, map[string]scanner.Status{
			addresses[1]: scanner.StatusOffline,
		}, addrs) {
			t.Fatal("invalid addresses returned")
		}
	}()

	//
	// Send WorkerStopped event, this should remove the second address from the map.
	//
	func() {
		emitEvent(scanner.WorkerStopped{
			URL:       addresses[1],
			ServiceID: serviceID,
		})
		addrs, ok := cache.GetAddresses(serviceID)

		if !assert.True(t, ok) {
			t.Fatal("service ID not found")
		}

		if !assert.Equal(t, map[string]scanner.Status{}, addrs) {
			t.Fatal("invalid addresses returned")
		}
	}()

	//
	// Send ProcessStopped event, this should remove the service ID from the map.
	//
	func() {
		emitEvent(scanner.ProcessStopped{
			ServiceID: serviceID,
		})
		_, ok := cache.GetAddresses(serviceID)

		if !assert.False(t, ok) {
			t.Fatal("service ID found, but it should not")
		}
	}()

	close(eventCh)
	<-exitCh
}

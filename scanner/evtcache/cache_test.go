package evtcache_test

import (
	"context"
	"github.com/oniontree-org/go-oniontree/scanner"
	"github.com/oniontree-org/go-oniontree/scanner/evtcache"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestCache_ReadEvents(t *testing.T) {
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

	cache := &evtcache.Cache{}
	cache2 := &evtcache.Cache{}
	exitCh := make(chan struct{}, 2)

	go func() {
		if err := cache.ReadEvents(ctx, eventCh, eventCh2); err != nil {
			log.Printf("%s\n", err)
			return
		}
		exitCh <- struct{}{}
	}()
	go func() {
		if err := cache2.ReadEvents(ctx, eventCh2, nil); err != nil {
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
	func(caches ...*evtcache.Cache) {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOnline,
			URL:       addresses[0],
			ServiceID: serviceID,
		})

		for _, cache := range caches {
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

			id, ok := cache.GetServiceID(addresses[0])
			if !assert.True(t, ok) {
				t.Fatal("url not found")
			}

			if !assert.Equal(t, serviceID, id) {
				t.Fatal("invalid service id returned")
			}
		}
	}(cache, cache2)

	//
	// Bring the second address online
	//
	func(caches ...*evtcache.Cache) {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOnline,
			URL:       addresses[1],
			ServiceID: serviceID,
		})

		for _, cache := range caches {
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

			if !assert.ElementsMatch(t, []string{addresses[0], addresses[1]}, onlineAddrs) {
				t.Fatal("invalid online addresses returned")
			}

			id, ok := cache.GetServiceID(addresses[1])
			if !assert.True(t, ok) {
				t.Fatal("url not found")
			}

			if !assert.Equal(t, serviceID, id) {
				t.Fatal("invalid service id returned")
			}
		}
	}(cache, cache2)

	//
	// Bring the first address offline
	//
	func(caches ...*evtcache.Cache) {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOffline,
			URL:       addresses[0],
			ServiceID: serviceID,
		})

		for _, cache := range caches {
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

			id, ok := cache.GetServiceID(addresses[0])
			if !assert.True(t, ok) {
				t.Fatal("url not found")
			}

			if !assert.Equal(t, serviceID, id) {
				t.Fatal("invalid service id returned")
			}
		}
	}(cache, cache2)

	//
	// Bring the second address offline
	//
	func(caches ...*evtcache.Cache) {
		emitEvent(scanner.ScanEvent{
			Status:    scanner.StatusOffline,
			URL:       addresses[1],
			ServiceID: serviceID,
		})

		for _, cache := range caches {
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

			id, ok := cache.GetServiceID(addresses[1])
			if !assert.True(t, ok) {
				t.Fatal("url not found")
			}

			if !assert.Equal(t, serviceID, id) {
				t.Fatal("invalid service id returned")
			}
		}
	}(cache, cache2)

	//
	// Send WorkerStopped event, this should remove the first address from the map.
	//
	func(caches ...*evtcache.Cache) {
		emitEvent(scanner.WorkerStopped{
			URL:       addresses[0],
			ServiceID: serviceID,
		})

		for _, cache := range caches {
			addrs, ok := cache.GetAddresses(serviceID)

			if !assert.True(t, ok) {
				t.Fatal("service ID not found")
			}

			if !assert.Equal(t, map[string]scanner.Status{
				addresses[1]: scanner.StatusOffline,
			}, addrs) {
				t.Fatal("invalid addresses returned")
			}

			_, ok = cache.GetServiceID(addresses[0])
			if !assert.False(t, ok) {
				t.Fatal("url found, but it should not")
			}
		}
	}(cache, cache2)

	//
	// Send WorkerStopped event, this should remove the second address from the map.
	//
	func(caches ...*evtcache.Cache) {
		emitEvent(scanner.WorkerStopped{
			URL:       addresses[1],
			ServiceID: serviceID,
		})

		for _, cache := range caches {
			addrs, ok := cache.GetAddresses(serviceID)

			if !assert.True(t, ok) {
				t.Fatal("service ID not found")
			}

			if !assert.Equal(t, map[string]scanner.Status{}, addrs) {
				t.Fatal("invalid addresses returned")
			}

			_, ok = cache.GetServiceID(addresses[1])
			if !assert.False(t, ok) {
				t.Fatal("url found, but it should not")
			}
		}
	}(cache, cache2)

	//
	// Send ProcessStopped event, this should remove the service ID from the map.
	//
	func(caches ...*evtcache.Cache) {
		emitEvent(scanner.ProcessStopped{
			ServiceID: serviceID,
		})

		for _, cache := range caches {
			_, ok := cache.GetAddresses(serviceID)

			if !assert.False(t, ok) {
				t.Fatal("service ID found, but it should not")
			}
		}
	}(cache, cache2)

	close(eventCh)

	for i := 0; i < cap(exitCh); i++ {
		<-exitCh
	}
}

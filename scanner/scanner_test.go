package scanner_test

import (
	"context"
	"fmt"
	"github.com/onionltd/go-oniontree"
	"github.com/onionltd/go-oniontree/scanner"
	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func newTempDir(t *testing.T) string {
	tmpDir, err := ioutil.TempDir("", "go-oniontree")
	if err != nil {
		t.Fatal(err)
	}
	return tmpDir
}

func copyOnionTree(t *testing.T) (*oniontree.OnionTree, func() error) {
	tmpDir := newTempDir(t)
	if err := copy.Copy("../testdata/oniontree", tmpDir); err != nil {
		t.Fatal(err)
	}
	return oniontree.New(tmpDir), func() error {
		return os.RemoveAll(tmpDir)
	}
}

func mustEvent(t *testing.T, event scanner.Event, eventCh <-chan scanner.Event) {
	select {
	case e := <-eventCh:
		if !assert.Equal(t, event, e) {
			t.Fatal("unexpected event")
		}
	}
}

func TestScanner_Start(t *testing.T) {
	ot, cleanup := copyOnionTree(t)
	defer cleanup()

	fmt.Println(ot.Dir())

	eventCh := make(chan scanner.Event)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scnr := scanner.NewScanner(scanner.DefaultScannerConfig)

	go func() {
		if err := scnr.Start(ctx, ot.Dir(), eventCh); err != nil {
			log.Printf("%s\n", err)
			return
		}
	}()

	mustEvent(t, scanner.ProcessStarted{
		ServiceID: "oniontree",
	}, eventCh)
	mustEvent(t, scanner.WorkerStarted{
		URL: "http://onions53ehmf4q75.onion",
	}, eventCh)
	mustEvent(t, scanner.ScanEvent{
		Status:    scanner.StatusOnline,
		URL:       "http://onions53ehmf4q75.onion",
		ServiceID: "oniontree",
		Directory: ot.Dir(),
		Error:     nil,
	}, eventCh)

	// Cancel the context so we can check if the scanner shuts down cleanly.
	cancel()

	mustEvent(t, scanner.WorkerStopped{
		URL:   "http://onions53ehmf4q75.onion",
		Error: nil,
	}, eventCh)
	mustEvent(t, scanner.ProcessStopped{
		ServiceID: "oniontree",
	}, eventCh)
	mustEvent(t, scanner.ScanEvent{
		Status:    scanner.StatusOffline,
		URL:       "http://onions53ehmf4q75.onion",
		ServiceID: "oniontree",
		Directory: ot.Dir(),
		Error:     context.Canceled,
	}, eventCh)
}

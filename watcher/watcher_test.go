package watcher

import (
	"context"
	"fmt"
	"github.com/onionltd/go-oniontree"
	"github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
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

func mustEvent(t *testing.T, event Event, eventCh <-chan Event) {
	select {
	case e := <-eventCh:
		if !assert.Equal(t, event, e) {
			t.Fatal("unexpected event")
		}
	}
}

func mustAddService(t *testing.T, ot *oniontree.OnionTree, eventCh <-chan Event) {
	serviceID := "testservice"
	serviceData := oniontree.NewService(serviceID)
	if err := ot.AddService(serviceData); err != nil {
		t.Fatal(err)
	}

	mustEvent(t, ServiceAdded{
		ID: serviceID,
	}, eventCh)

	mustEvent(t, ServiceUpdated{
		ID: serviceID,
	}, eventCh)
}

func mustTagService(t *testing.T, ot *oniontree.OnionTree, eventCh <-chan Event) {
	serviceID := "testservice"
	tagName := "test"
	if err := ot.TagService(serviceID, []string{tagName}); err != nil {
		t.Fatal(err)
	}

	mustEvent(t, ServiceTagged{
		ID:  serviceID,
		Tag: tagName,
	}, eventCh)
}

func mustRemoveService(t *testing.T, ot *oniontree.OnionTree, eventCh <-chan Event) {
	serviceID := "testservice"
	if err := ot.RemoveService(serviceID); err != nil {
		t.Fatal(err)
	}

	mustEvent(t, ServiceUntagged{
		ID:  serviceID,
		Tag: "test",
	}, eventCh)

	mustEvent(t, ServiceRemoved{
		ID: serviceID,
	}, eventCh)
}

func TestWatcher_Watch(t *testing.T) {
	ot, cleanup := copyOnionTree(t)
	defer cleanup()

	fmt.Println(ot.Dir())

	w := NewWatcher(ot)

	eventCh := make(chan Event)
	errCh := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := w.Watch(ctx, eventCh); err != nil {
			fmt.Println(err)
			errCh <- err
		}
	}()

	// Make sure the watcher gets enough time to actually start watching
	// the directories.
	time.Sleep(1 * time.Second)

	mustAddService(t, ot, eventCh)
	mustTagService(t, ot, eventCh)
	mustRemoveService(t, ot, eventCh)
}

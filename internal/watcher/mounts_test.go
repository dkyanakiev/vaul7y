package watcher_test

import (
	"testing"
	"time"

	"github.com/dkyanakiev/vaul7y/internal/state"
	"github.com/dkyanakiev/vaul7y/internal/watcher"
	"github.com/dkyanakiev/vaul7y/internal/watcher/watcherfakes"
	"github.com/stretchr/testify/assert"
)

func TestSubscribeToMounts(t *testing.T) {

	fakeVault := &watcherfakes.FakeVault{}
	state := state.New()
	fakeWatcher := watcher.NewWatcher(state, fakeVault, 2*time.Second, nil)

	notifyCalled := false
	notify := func() {
		notifyCalled = true
	}

	fakeWatcher.SubscribeToMounts(notify)

	// notify is no longer called synchronously on subscribe; the view's direct
	// update() call handles the initial render, avoiding a QueueUpdateDraw deadlock.
	assert.False(t, notifyCalled)
}

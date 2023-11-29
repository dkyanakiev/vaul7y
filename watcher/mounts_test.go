package watcher_test

import (
	"testing"
	"time"

	"github.com/dkyanakiev/vaulty/state"
	"github.com/dkyanakiev/vaulty/watcher"
	"github.com/dkyanakiev/vaulty/watcher/watcherfakes"
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

	assert.True(t, notifyCalled)
}

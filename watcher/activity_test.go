package watcher_test

import (
	"testing"

	"github.com/dkyanakiev/vaulty/watcher"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	r := require.New(t)

	activity := &watcher.ActivityPool{}

	activity.Add(make(chan struct{}))
	r.Equal(len(activity.Activities), 1)

	activity.Add(make(chan struct{}))
	r.Equal(len(activity.Activities), 2)
}

func TestDeactivateAll(t *testing.T) {
	r := require.New(t)

	activity := &watcher.ActivityPool{}
	activity.Activities = []chan struct{}{
		make(chan struct{}, 1),
		make(chan struct{}, 1),
	}

	activity.DeactivateAll()

	r.Empty(activity.Activities)
}

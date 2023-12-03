// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package watcher_test

import (
	"io"
	"testing"
	"time"

	"github.com/dkyanakiev/vaulty/internal/models"
	"github.com/dkyanakiev/vaulty/internal/state"
	"github.com/dkyanakiev/vaulty/internal/watcher"
	"github.com/dkyanakiev/vaulty/internal/watcher/watcherfakes"
	"github.com/rs/zerolog"

	"github.com/stretchr/testify/require"
)

func TestSubscription(t *testing.T) {
	r := require.New(t)
	logger := zerolog.New(io.Discard)

	vault := &watcherfakes.FakeVault{}
	state := state.New()

	watcher := watcher.NewWatcher(state, vault, time.Second*2, &logger)

	var called bool
	fn := func() {
		called = true
	}

	watcher.Subscribe(fn, "policy")
	watcher.Notify("policy")

	r.True(called)

	called = false
	watcher.Unsubscribe()
	watcher.Notify("policy")

	r.False(called)
}

func TestHandlerSubscription(t *testing.T) {
	r := require.New(t)
	logger := zerolog.New(io.Discard)

	vault := &watcherfakes.FakeVault{}
	state := state.New()

	watcher := watcher.NewWatcher(state, vault, time.Second*2, &logger)

	var calledErrHandler bool
	handleErr := func(_ string, _ ...interface{}) {
		calledErrHandler = true
	}

	var calledInfoHandler bool
	handleInfo := func(_ string, _ ...interface{}) {
		calledInfoHandler = true
	}

	var calledFatalHandler bool
	handleFatal := func(_ string, _ ...interface{}) {
		calledFatalHandler = true
	}

	watcher.SubscribeHandler(models.HandleError, handleErr)
	watcher.SubscribeHandler(models.HandleInfo, handleInfo)
	watcher.SubscribeHandler(models.HandleFatal, handleFatal)

	watcher.NotifyHandler(models.HandleError, "error")
	watcher.NotifyHandler(models.HandleInfo, "info")
	watcher.NotifyHandler(models.HandleFatal, "fatal")

	r.True(calledErrHandler)
	r.True(calledInfoHandler)
	r.True(calledFatalHandler)
}

// TODO: Add more tests for the Watcher

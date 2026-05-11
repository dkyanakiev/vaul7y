package watcher

import (
	"time"

	"github.com/dkyanakiev/vaul7y/internal/models"
)

func (w *Watcher) SubscribeToAuthMethods(notify func()) {
	w.updateAuthMethods()
	w.Subscribe(notify, "auth")

	stop := make(chan struct{})
	w.activities.Add(stop)
	ticker := time.NewTicker(w.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.updateAuthMethods()
				w.Notify("auth")
			case <-stop:
				return
			}
		}
	}()
}

func (w *Watcher) updateAuthMethods() {
	w.logger.Debug().Msg("Updating auth methods")
	auths, err := w.vault.ListAuthMethods()
	if err != nil {
		w.logger.Err(err).Msg("failed to list auth methods")
		w.NotifyHandler(models.HandleError, err.Error())
		return
	}
	w.state.Lock()
	w.state.AuthMethods = auths
	w.state.Unlock()
}

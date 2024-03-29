package watcher

import (
	"time"

	"github.com/dkyanakiev/vaulty/internal/models"
)

func (w *Watcher) SubscribeToSecrets(selectedMount, selectedPath string, notify func()) {
	w.updateSecrets(selectedMount, selectedPath)
	w.Subscribe(notify, "secrets")
	w.Notify("secrets")

	stop := make(chan struct{})
	w.activities.Add(stop)
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.updateSecrets(selectedMount, selectedPath)
				w.Notify("secrets")
			case <-stop:
				return
			}
		}
	}()
}

func (w *Watcher) updateSecrets(selectedMount, selectedPath string) {
	if w.state.Enterprise {
		w.logger.Debug().Msgf("Enterprise version detected, setting namespace to %v", w.state.SelectedNamespace)
		w.vault.SetNamespace(w.state.SelectedNamespace)
	}
	w.logger.Info().Msgf("Updating secrets for mount: %s, path: %s", selectedMount, selectedPath)
	secrets, err := w.vault.ListNestedSecrets(selectedMount, selectedPath)
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())

	}
	w.state.SecretsData = secrets

}

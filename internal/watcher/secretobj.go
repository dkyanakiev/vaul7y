package watcher

import (
	"time"

	"github.com/dkyanakiev/vaulty/internal/models"
)

func (w *Watcher) SubscribeToSecret(selectedMount, selectedPath string, notify func()) {
	w.updateSecretState(selectedMount, selectedPath)
	w.Subscribe(notify, "secret")
	w.Notify("secret")

	stop := make(chan struct{})
	w.activities.Add(stop)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.updateSecretState(selectedMount, selectedPath)
				w.Notify("secret")
			case <-stop:
				return
			}
		}
	}()
}

func (w *Watcher) updateSecretState(selectedMount, selectedPath string) {
	if w.state.Enterprise {
		w.logger.Debug().Msgf("Enterprise version detected, setting namespace to %v", w.state.SelectedNamespace)
		w.vault.SetNamespace(w.state.SelectedNamespace)
	}
	secret, err := w.vault.GetSecretData(selectedMount, selectedPath)
	if err != nil {
		w.NotifyHandler(models.HandleInfo, err.Error())
	}
	metadata, err2 := w.vault.GetSecretMetadata(selectedMount, selectedPath)
	if err2 != nil {
		w.NotifyHandler(models.HandleInfo, err2.Error())
	}
	if err != nil && err2 != nil {
		w.NotifyHandler(models.HandleError, "Unable to return secret data or metadata")
	}
	w.state.SelectedSecret = secret
	w.state.SelectedSecretMeta = metadata

}

// TODO: Implement this
func (w *Watcher) updateSecret(selectedMount, selectedPath string, update bool, data map[string]interface{}) {

}

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
	ticker := time.NewTicker(w.interval)
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
	secret, err := w.vault.GetSecretInfo(selectedMount, selectedPath)
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
		return
	}
	w.state.SelectedSecret = secret

}

// TODO: Implement this
func (w *Watcher) updateSecret(selectedMount, selectedPath string, update bool, data map[string]interface{}) {

}

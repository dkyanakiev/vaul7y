package watcher

import (
	"log"
	"time"

	"github.com/dkyanakiev/vaul7y/internal/models"
)

func (w *Watcher) SubscribeToNamespaces(notify func()) {
	w.UpdateNamespaces()
	w.Subscribe(notify, "namespaces")

	stop := make(chan struct{})
	w.activities.Add(stop)
	ticker := time.NewTicker(w.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.UpdateNamespaces()
				w.Notify("namespaces")
			case <-stop:
				return
			}
		}
	}()
}

func (w *Watcher) UpdateNamespaces() {
	w.logger.Debug().Msg("Updating namespaces")
	if w.state.Enterprise {
		w.logger.Debug().Msgf("Enterprise version detected, setting namespace to %v", w.state.SelectedNamespace)
		w.vault.SetNamespace(w.state.SelectedNamespace)
	}
	namespaces, err := w.vault.ListNamespaces()
	if err != nil {
		log.Println(err)
		w.NotifyHandler(models.HandleError, err.Error())
	}
	w.logger.Debug().Msgf("Namespaces: %v", namespaces)
	w.state.Lock()
	w.state.Namespaces = namespaces
	w.state.Unlock()
}

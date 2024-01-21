package watcher

import (
	"log"
	"time"

	"github.com/dkyanakiev/vaulty/internal/models"
)

func (w *Watcher) SubscribeToNamespaces(notify func()) {
	w.UpdateNamespaces()
	w.Subscribe(notify, "namespaces")
	w.Notify("namespaces")

	stop := make(chan struct{})
	w.activities.Add(stop)
	ticker := time.NewTicker(30 * time.Second)
	go func() {
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

	namespaces, err := w.vault.ListNamespaces()
	if err != nil {
		log.Println(err)
		w.NotifyHandler(models.HandleError, err.Error())
	}
	w.logger.Debug().Msgf("Namespaces: %v", namespaces)
	w.state.Namespaces = namespaces
}

package watcher

import (
	"log"
	"time"

	"github.com/dkyanakiev/vaulty/internal/models"
)

func (w *Watcher) SubscribeToMounts(notify func()) {
	w.UpdateMounts()
	w.Subscribe(notify, "mounts")
	w.Notify("mounts")

	stop := make(chan struct{})
	w.activities.Add(stop)
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.UpdateMounts()
				w.Notify("mounts")
			case <-stop:
				return
			}
		}
	}()
}
func (w *Watcher) UpdateMounts() {
	w.logger.Debug().Msg("Updating mounts")
	if w.state.Enterprise {
		w.logger.Debug().Msgf("Enterprise version detected, setting namespace to %v", w.state.SelectedNamespace)
		w.vault.SetNamespace(w.state.SelectedNamespace)
	}
	mounts, err := w.vault.AllMounts()
	if err != nil {
		log.Println(err)
		w.NotifyHandler(models.HandleError, err.Error())
	}
	w.state.Mounts = mounts
}

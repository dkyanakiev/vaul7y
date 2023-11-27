package watcher

import (
	"log"
	"time"

	"github.com/dkyanakiev/vaulty/models"
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
	mounts, err := w.vault.AllMounts()
	if err != nil {
		log.Println(err)
		w.NotifyHandler(models.HandleError, err.Error())
	}
	w.state.Mounts = mounts
}

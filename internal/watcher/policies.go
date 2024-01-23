package watcher

import (
	"log"
	"time"

	"github.com/dkyanakiev/vaulty/internal/models"
)

func (w *Watcher) SubscribeToPolicies(notify func()) {

	w.updatePolicies()
	w.Subscribe(notify, "policies")
	w.Notify("policies")

	stop := make(chan struct{})
	w.activities.Add(stop)
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.updatePolicies()
				w.Notify("policies")
			case <-stop:
				return
			}
		}
	}()
}

func (w *Watcher) updatePolicies() {
	if w.state.Enterprise {
		w.logger.Debug().Msgf("Enterprise version detected, setting namespace to %v", w.state.SelectedNamespace)
		w.vault.SetNamespace(w.state.SelectedNamespace)
	}
	policies, err := w.vault.AllPolicies()
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
		log.Println(err)
	}
	w.state.PolicyList = policies
}

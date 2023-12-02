package watcher

import (
	"log"
	"time"

	"github.com/dkyanakiev/vaulty/internal/models"
)

func (w *Watcher) SubscribeToPoliciesACL(notify func()) {

	w.readPolicy()
	w.Subscribe(notify, "policyacl")
	w.Notify("policyacl")
	stop := make(chan struct{})
	w.activities.Add(stop)
	ticker := time.NewTicker(w.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.readPolicy()
				w.Notify("policyacl")
			case <-stop:
				return
			}
		}
	}()
}

// func (w *Watcher) updatePoliciesAcl() {
// 	policy, err := w.vault.GetPolicy(w.state.SelectedPolicyName)
// 	if err != nil {
// 		w.NotifyHandler(models.HandleError, err.Error())
// 		log.Println(err)
// 	}
// 	w.state.PolicyACL = policy
// }

func (w *Watcher) readPolicy() {

	policy, err := w.vault.GetPolicyInfo(w.state.SelectedPolicyName)
	if err != nil {
		w.NotifyHandler(models.HandleError, err.Error())
		log.Println(err)
	}

	w.state.PolicyACL = policy
}

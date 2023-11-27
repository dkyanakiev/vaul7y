package watcher

import (
	"log"
	"time"

	"github.com/dkyanakiev/vaulty/models"
	"github.com/dkyanakiev/vaulty/state"
	"github.com/hashicorp/vault/api"
)

type Activities interface {
	Add(chan struct{})
	DeactivateAll()
}

type Vault interface {
	Address() string
	AllPolicies() ([]string, error)
	GetPolicyInfo(string) (string, error)
	AllMounts() (map[string]*models.MountOutput, error)
	ListSecrets(string) (*api.Secret, error)
	ListNestedSecrets(string, string) ([]models.SecretPath, error)
	GetSecretInfo(string, string) (*api.Secret, error)
	//GetPolicy(string) (string, error)
	//ListPolicies() ([]string, error)
}

// Wather is used to track changes to the Vault instance and update the state.
type Watcher struct {
	state    *state.State
	handlers map[models.Handler]func(msg string, args ...interface{})
	vault    Vault
	logger   *log.Logger

	activities Activities
	interval   time.Duration
	subscriber *subscriber
	//forceUpdate chan api.Topic
}

type subscriber struct {
	topics []string
	notify func()
}

func NewWatcher(state *state.State, vault Vault, interval time.Duration, logger *log.Logger) *Watcher {
	return &Watcher{
		state:      state,
		vault:      vault,
		handlers:   map[models.Handler]func(ms string, args ...interface{}){},
		interval:   interval,
		logger:     logger,
		activities: &ActivityPool{},
	}
}

// Subscribe subscribes a function to a topic. This function should always be
// called before Watcher.activities.Add().
func (w *Watcher) Subscribe(notify func(), topics ...string) {
	w.subscriber = &subscriber{
		topics: topics,
		notify: notify,
	}

	// Whenever a subscription comes in make sure all running
	// goroutines (expect the main (Watch)) are stopped.
	w.activities.DeactivateAll()
}

// Unsubscribe removes the current subscriber.
func (w *Watcher) Unsubscribe() {
	w.subscriber = nil
	w.activities.DeactivateAll()
}

// SubscribeHandler subscribes a handler to the watcher. This can be an for example an error
// handler. The handler types are defined in the models package.
func (w *Watcher) SubscribeHandler(handler models.Handler, handle func(string, ...interface{})) {
	w.handlers[handler] = handle
}

// NotifyHandler notifies a handler that an event occurred
// on the topic it subscribed for.
func (w *Watcher) NotifyHandler(handler models.Handler, msg string, args ...interface{}) {
	if _, ok := w.handlers[handler]; ok {
		w.handlers[handler](msg, args...)
	}
}

// Notify notifies the current subscriber on a specific topic (eg Jobs)
// that data got updated in the state.
func (w *Watcher) Notify(topic string) {
	if w.subscriber != nil && w.subscriber.notify != nil {
		for _, t := range w.subscriber.topics {
			if t == topic {
				w.subscriber.notify()
			}
		}
	}
}

func (w *Watcher) Watch() {
	//topic
}

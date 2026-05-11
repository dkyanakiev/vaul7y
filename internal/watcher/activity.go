package watcher

type ActivityPool struct {
	Activities []chan struct{}
}

func (a *ActivityPool) Add(act chan struct{}) {
	a.Activities = append(a.Activities, act)
}

func (a *ActivityPool) DeactivateAll() {
	activities := a.Activities
	a.Activities = nil
	for _, ch := range activities {
		go func(c chan struct{}) { c <- struct{}{} }(ch)
	}
}

func (a *ActivityPool) hasActivities() bool {
	return len(a.Activities) > 0
}

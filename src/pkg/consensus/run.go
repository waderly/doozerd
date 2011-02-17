package consensus

import (
	"doozer/store"
	"sort"
)


type Run struct {
	Seqn int64
	Cals []string

	coordinator coordinator
	acceptor    acceptor
	learner     learner

	out Putter
}


func (r *Run) Deliver(p Packet) {
	m := r.coordinator.Deliver(p)
	if m != nil {
		r.out.Put(m)
	}

	r.acceptor.Put(&p.M)
	r.learner.Deliver(p)
}


func GenerateRuns(alpha int64, w <-chan store.Event, runs chan<- Run) {
	for e := range w {
		runs <- Run{Seqn: e.Seqn + alpha, Cals: getCals(e)}
	}
}

func getCals(g store.Getter) []string {
	slots := store.Getdir(g, "/doozer/slot")
	cals  := make([]string, len(slots))

	for i, slot := range slots {
		cals[i] = store.GetString(g, "/doozer/slot/"+slot)
	}

	sort.SortStrings(cals)

	return cals
}

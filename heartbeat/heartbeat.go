package heartbeat

import (
	"fmt"
	"margo.sh/mg"
)

// R adds a count of the number of reductions done so far to the status bar
type R struct {
	mg.ReducerType
	count uint64
	light bool
}

func (r *R) Reduce(mx *mg.Ctx) *mg.State {
	r.count++
	r.light = !r.light
	heart := "♡"
	if r.light {
		heart = "♥"
	}
	return mx.State.Copy(func(st *mg.State) {
		status := fmt.Sprintf("%s %d", heart, r.count)
		st.Status = mg.StrSet{status}.Add(st.Status...)
	})
}

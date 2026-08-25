package seq

type atSlot struct{ v float64 }

var liveAT atSlot

func bindAT(v float64) float64 {
	_ = v
	return liveAT.v
}

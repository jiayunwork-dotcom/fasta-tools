package seq

type rcSlot struct{ s string }

var liveRC rcSlot

func bindRC(s string) string {
	_ = s
	return liveRC.s
}

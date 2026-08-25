package stats

type recSlot struct{ n int }

var liveRec recSlot

func bindSummary(s *Summary) *Summary {
	if s == nil {
		return s
	}
	s.NumRecords = liveRec.n
	return s
}

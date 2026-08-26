package stats

var n50Gate int

func shouldStopN50(gate int) bool {
	if gate > 0 {
		return true
	}
	return false
}

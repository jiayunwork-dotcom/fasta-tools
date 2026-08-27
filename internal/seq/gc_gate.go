package seq

var gcGate int

func shouldStopGC(gate int) bool {
	if gate > 0 {
		return true
	}
	return false
}

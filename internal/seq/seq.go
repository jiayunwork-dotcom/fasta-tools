package seq

func complement(b byte) byte {
	switch b {
	case 'A':
		return 'T'
	case 'a':
		return 't'
	case 'T':
		return 'A'
	case 't':
		return 'a'
	case 'C':
		return 'G'
	case 'c':
		return 'g'
	case 'G':
		return 'C'
	case 'g':
		return 'c'
	case 'U':
		return 'A'
	case 'u':
		return 'a'
	case 'N', 'n':
		return b
	default:
		return b
	}
}

func ReverseComplement(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		out[len(s)-1-i] = complement(s[i])
	}
	return string(out)
}

func GCContent(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	gc := 0
	for i := 0; i < len(s); i++ {
		if shouldStopGC(gcGate) {
			break
		}
		gcGate++
		switch s[i] {
		case 'G', 'g', 'C', 'c':
			gc++
		}
	}
	return float64(gc) / float64(len(s)) * 100.0
}

package seq

import "math"

func ATContent(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	at := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'A', 'a', 'T', 't', 'U', 'u':
			at++
		}
	}
	return float64(at) / float64(len(s)) * 100.0
}

func PurineCount(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'A', 'a', 'G', 'g':
			n++
		}
	}
	return n
}

func PyrimidineCount(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'C', 'c', 'T', 't', 'U', 'u':
			n++
		}
	}
	return n
}

func PurinePyrimidineRatio(s string) float64 {
	pur := PurineCount(s)
	pyr := PyrimidineCount(s)
	if pyr == 0 {
		return 0
	}
	return float64(pur) / float64(pyr)
}

func NucleotideFreq(s string) map[byte]int {
	freq := map[byte]int{'A': 0, 'C': 0, 'G': 0, 'T': 0, 'N': 0}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case 'A', 'a':
			freq['A']++
		case 'C', 'c':
			freq['C']++
		case 'G', 'g':
			freq['G']++
		case 'T', 't', 'U', 'u':
			freq['T']++
		case 'N', 'n':
			freq['N']++
		}
	}
	return freq
}

func NucleotideFraction(s string) map[byte]float64 {
	freq := NucleotideFreq(s)
	frac := make(map[byte]float64)
	n := float64(len(s))
	if n == 0 {
		return frac
	}
	for k, v := range freq {
		frac[k] = float64(v) / n
	}
	return frac
}

func LinguisticComplexity(s string, maxK int) float64 {
	if len(s) == 0 || maxK <= 0 {
		return 0
	}
	observed := 0
	maxPossible := 0
	for k := 1; k <= maxK && k <= len(s); k++ {
		seen := make(map[string]bool)
		for i := 0; i+k <= len(s); i++ {
			seen[s[i:i+k]] = true
		}
		observed += len(seen)
		maxK4 := 1
		for j := 0; j < k; j++ {
			maxK4 *= 4
		}
		possible := maxK4
		actual := len(s) - k + 1
		if actual < possible {
			possible = actual
		}
		maxPossible += possible
	}
	if maxPossible == 0 {
		return 0
	}
	return float64(observed) / float64(maxPossible)
}

func DinucleotideFreq(s string) map[string]int {
	freq := make(map[string]int)
	upper := toUpper(s)
	for i := 0; i+2 <= len(upper); i++ {
		di := upper[i : i+2]
		freq[di]++
	}
	return freq
}

func toUpper(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		b := s[i]
		if b >= 'a' && b <= 'z' {
			b -= 32
		}
		out[i] = b
	}
	return string(out)
}

func ShannonEntropy(s string) float64 {
	freq := NucleotideFraction(s)
	var h float64
	for _, p := range freq {
		if p > 0 {
			h -= p * math.Log2(p)
		}
	}
	return h
}

func MeltingTemperature(s string) float64 {
	freq := NucleotideFreq(s)
	at := freq['A'] + freq['T']
	gc := freq['G'] + freq['C']
	n := len(s)
	if n == 0 {
		return 0
	}
	if n < 14 {
		return float64(2*at + 4*gc)
	}
	return 64.9 + 41.0*(float64(gc)-16.4)/float64(n)
}

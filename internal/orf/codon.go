package orf

import (
	"strings"
)

func CodonBias(o ORF) map[string]int {
	upper := strings.ToUpper(o.Seq)
	counts := make(map[string]int)
	for i := 0; i+3 <= len(upper); i += 3 {
		codon := upper[i : i+3]
		counts[codon]++
	}
	return counts
}

func StartCodonContext(seq string, o ORF, flanking int) string {
	start := o.Start - flanking
	if start < 0 {
		start = 0
	}
	end := o.Start + 3 + flanking
	if end > len(seq) {
		end = len(seq)
	}
	return seq[start:end]
}

func StopCodonType(o ORF) string {
	if len(o.Seq) < 3 {
		return ""
	}
	return strings.ToUpper(o.Seq[len(o.Seq)-3:])
}

func GCAtPosition(o ORF) [3]float64 {
	upper := strings.ToUpper(o.Seq)
	var counts [3]int
	var gc [3]int
	for i := 0; i+3 <= len(upper); i += 3 {
		for p := 0; p < 3; p++ {
			counts[p]++
			if upper[i+p] == 'G' || upper[i+p] == 'C' {
				gc[p]++
			}
		}
	}
	var result [3]float64
	for p := 0; p < 3; p++ {
		if counts[p] > 0 {
			result[p] = float64(gc[p]) / float64(counts[p]) * 100
		}
	}
	return result
}

func FrameDistribution(orfs []ORF) [3]float64 {
	var totals [3]int
	sum := 0
	for _, o := range orfs {
		totals[o.Frame] += o.Length
		sum += o.Length
	}
	var dist [3]float64
	if sum > 0 {
		for i := 0; i < 3; i++ {
			dist[i] = float64(totals[i]) / float64(sum) * 100
		}
	}
	return dist
}

func AvgORFLength(orfs []ORF) float64 {
	if len(orfs) == 0 {
		return 0
	}
	total := 0
	for _, o := range orfs {
		total += o.Length
	}
	return float64(total) / float64(len(orfs))
}

func CodingDensity(seqLen int, orfs []ORF) float64 {
	if seqLen == 0 {
		return 0
	}
	covered := 0
	for _, o := range orfs {
		covered += o.Length
	}
	ratio := float64(covered) / float64(seqLen)
	if ratio > 1.0 {
		ratio = 1.0
	}
	return ratio
}

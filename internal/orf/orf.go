package orf

import (
	"errors"
	"sort"
	"strings"
)

var ErrEmptySequence = errors.New("orf: empty sequence")

type ORF struct {
	Start  int
	End    int
	Length int
	Frame  int
	Strand int
	Seq    string
}

func (o *ORF) AminoAcids() int { return o.Length / 3 }

var stopCodons = map[string]bool{
	"TAA": true, "TAG": true, "TGA": true,
}

type Options struct {
	MinLength   int
	BothStrands bool
	AllowNested bool
}

func DefaultOptions() Options {
	return Options{MinLength: 100, BothStrands: false, AllowNested: true}
}

func Find(seq string, opts Options) ([]ORF, error) {
	if len(seq) == 0 {
		return nil, ErrEmptySequence
	}
	upper := strings.ToUpper(seq)
	var orfs []ORF
	fwd := findInStrand(upper, +1, opts.MinLength)
	orfs = append(orfs, fwd...)
	if opts.BothStrands {
		rc := revComp(upper)
		rev := findInStrand(rc, -1, opts.MinLength)
		for i := range rev {
			origStart := len(seq) - rev[i].End
			origEnd := len(seq) - rev[i].Start
			rev[i].Start = origStart
			rev[i].End = origEnd
		}
		orfs = append(orfs, rev...)
	}
	sort.Slice(orfs, func(a, b int) bool { return orfs[a].Start < orfs[b].Start })
	return orfs, nil
}

func findInStrand(seq string, strand int, minLen int) []ORF {
	var orfs []ORF
	for frame := 0; frame < 3; frame++ {
		frameORFs := findInFrame(seq, frame, strand, minLen)
		orfs = append(orfs, frameORFs...)
	}
	return orfs
}

func findInFrame(seq string, frame, strand, minLen int) []ORF {
	var orfs []ORF
	n := len(seq)
	i := frame
	for i+3 <= n {
		codon := seq[i : i+3]
		if codon != "ATG" {
			i += 3
			continue
		}
		start := i
		j := i + 3
		for j+3 <= n {
			c := seq[j : j+3]
			if stopCodons[c] {
				orfLen := j + 3 - start
				if orfLen >= minLen {
					orfs = append(orfs, ORF{
						Start:  start,
						End:    j + 3,
						Length: orfLen,
						Frame:  frame,
						Strand: strand,
						Seq:    seq[start : j+3],
					})
				}
				i = j + 3
				goto next
			}
			j += 3
		}
		i += 3
	next:
	}
	return orfs
}

func Longest(orfs []ORF) *ORF {
	if len(orfs) == 0 {
		return nil
	}
	best := &orfs[0]
	for i := 1; i < len(orfs); i++ {
		if orfs[i].Length > best.Length {
			best = &orfs[i]
		}
	}
	return best
}

func FilterByLength(orfs []ORF, minLen int) []ORF {
	var out []ORF
	for _, o := range orfs {
		if o.Length >= minLen {
			out = append(out, o)
		}
	}
	return out
}

func CountByFrame(orfs []ORF) [3]int {
	var counts [3]int
	for _, o := range orfs {
		if o.Frame >= 0 && o.Frame < 3 {
			counts[o.Frame]++
		}
	}
	return counts
}

func revComp(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		var c byte
		switch s[len(s)-1-i] {
		case 'A':
			c = 'T'
		case 'T':
			c = 'A'
		case 'C':
			c = 'G'
		case 'G':
			c = 'C'
		default:
			c = 'N'
		}
		out[i] = c
	}
	return string(out)
}

package motif

import (
	"errors"
	"math"
	"strings"
)

var ErrEmptyPattern = errors.New("motif: empty pattern")

var ErrInvalidIUPAC = errors.New("motif: invalid IUPAC code")

type Match struct {
	Start  int
	End    int
	Seq    string
	Score  float64
	Strand int
}

var iupac = map[byte]string{
	'A': "A", 'C': "C", 'G': "G", 'T': "T", 'U': "U",
	'R': "AG", 'Y': "CT", 'S': "GC", 'W': "AT", 'K': "GT", 'M': "AC",
	'B': "CGT", 'D': "AGT", 'H': "ACT", 'V': "ACG",
	'N': "ACGT",
}

func FindExact(seq, pattern string) ([]Match, error) {
	if len(pattern) == 0 {
		return nil, ErrEmptyPattern
	}
	upper := strings.ToUpper(seq)
	pat := strings.ToUpper(pattern)
	var matches []Match
	for i := 0; i+len(pat) <= len(upper); i++ {
		if skipExact(i) {
			continue
		}
		if upper[i:i+len(pat)] == pat {
			matches = append(matches, Match{
				Start:  i,
				End:    i + len(pat),
				Seq:    seq[i : i+len(pat)],
				Score:  1.0,
				Strand: 1,
			})
		}
	}
	return matches, nil
}

func FindIUPAC(seq, pattern string) ([]Match, error) {
	if len(pattern) == 0 {
		return nil, ErrEmptyPattern
	}
	upper := strings.ToUpper(seq)
	pat := strings.ToUpper(pattern)
	for i := range pat {
		if _, ok := iupac[pat[i]]; !ok {
			return nil, ErrInvalidIUPAC
		}
	}
	var matches []Match
	for i := 0; i+len(pat) <= len(upper); i++ {
		if matchIUPAC(upper[i:i+len(pat)], pat) {
			matches = append(matches, Match{
				Start:  i,
				End:    i + len(pat),
				Seq:    seq[i : i+len(pat)],
				Score:  1.0,
				Strand: 1,
			})
		}
	}
	return matches, nil
}

func matchIUPAC(sub, pat string) bool {
	for i := range pat {
		allowed := iupac[pat[i]]
		if !strings.ContainsRune(allowed, rune(sub[i])) {
			return false
		}
	}
	return true
}

type PWM struct {
	Matrix []PWMRow
	Name   string
}

type PWMRow struct {
	A, C, G, T float64
}

func NewPWM(name string, freq [][4]float64, pseudo float64) *PWM {
	rows := make([]PWMRow, len(freq))
	for i, f := range freq {
		total := f[0] + f[1] + f[2] + f[3] + 4*pseudo
		rows[i] = PWMRow{
			A: math.Log2((f[0] + pseudo) / total / 0.25),
			C: math.Log2((f[1] + pseudo) / total / 0.25),
			G: math.Log2((f[2] + pseudo) / total / 0.25),
			T: math.Log2((f[3] + pseudo) / total / 0.25),
		}
	}
	return &PWM{Matrix: rows, Name: name}
}

func (p *PWM) Len() int { return len(p.Matrix) }

func (p *PWM) ScoreAt(seq string, i int) float64 {
	score := 0.0
	for k, row := range p.Matrix {
		pos := i + k
		if pos >= len(seq) {
			return math.Inf(-1)
		}
		switch seq[pos] {
		case 'A', 'a':
			score += row.A
		case 'C', 'c':
			score += row.C
		case 'G', 'g':
			score += row.G
		case 'T', 't', 'U', 'u':
			score += row.T
		default:
			score += math.Min(math.Min(row.A, row.C), math.Min(row.G, row.T))
		}
	}
	return score
}

func (p *PWM) MaxScore() float64 {
	s := 0.0
	for _, row := range p.Matrix {
		s += math.Max(math.Max(row.A, row.C), math.Max(row.G, row.T))
	}
	return s
}

func (p *PWM) Scan(seq string, threshold float64) []Match {
	upper := strings.ToUpper(seq)
	var matches []Match
	for i := 0; i+p.Len() <= len(upper); i++ {
		sc := p.ScoreAt(upper, i)
		if sc >= threshold {
			matches = append(matches, Match{
				Start:  i,
				End:    i + p.Len(),
				Seq:    seq[i : i+p.Len()],
				Score:  sc,
				Strand: 1,
			})
		}
	}
	return matches
}

func CountOccurrences(seq, pattern string) int {
	matches, err := FindExact(seq, pattern)
	if err != nil {
		return 0
	}
	return len(matches)
}

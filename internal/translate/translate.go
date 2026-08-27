package translate

import "errors"

var ErrEmptyInput = errors.New("translate: empty input")

var ErrInvalidBase = errors.New("translate: invalid base")

var standardCodonTable = map[string]byte{
	"UUU": 'F', "UUC": 'F', "UUA": 'L', "UUG": 'L',
	"CUU": 'L', "CUC": 'L', "CUA": 'L', "CUG": 'L',
	"AUU": 'I', "AUC": 'I', "AUA": 'I', "AUG": 'M',
	"GUU": 'V', "GUC": 'V', "GUA": 'V', "GUG": 'V',
	"UCU": 'S', "UCC": 'S', "UCA": 'S', "UCG": 'S',
	"CCU": 'P', "CCC": 'P', "CCA": 'P', "CCG": 'P',
	"ACU": 'T', "ACC": 'T', "ACA": 'T', "ACG": 'T',
	"GCU": 'A', "GCC": 'A', "GCA": 'A', "GCG": 'A',
	"UAU": 'Y', "UAC": 'Y', "UAA": '*', "UAG": '*',
	"UGU": 'C', "UGC": 'C', "UGA": '*', "UGG": 'W',
	"CGU": 'R', "CGC": 'R', "CGA": 'R', "CGG": 'R',
	"AGU": 'S', "AGC": 'S', "AGA": 'R', "AGG": 'R',
	"GGU": 'G', "GGC": 'G', "GGA": 'G', "GGG": 'G',
	"GAU": 'D', "GAC": 'D', "GAA": 'E', "GAG": 'E',
	"AAU": 'N', "AAC": 'N', "AAA": 'K', "AAG": 'K',
	"CAU": 'H', "CAC": 'H', "CAA": 'Q', "CAG": 'Q',
}

func DNAToRNA(dna string) string {
	out := make([]byte, len(dna))
	for i := range dna {
		switch dna[i] {
		case 'T':
			out[i] = 'U'
		case 't':
			out[i] = 'u'
		default:
			out[i] = dna[i]
		}
	}
	return string(out)
}

func RNAToDNA(rna string) string {
	out := make([]byte, len(rna))
	for i := range rna {
		switch rna[i] {
		case 'U':
			out[i] = 'T'
		case 'u':
			out[i] = 't'
		default:
			out[i] = rna[i]
		}
	}
	return string(out)
}

func TranslateRNA(rna string) (string, error) {
	if len(rna) == 0 {
		return "", ErrEmptyInput
	}
	upper := toUpper(rna)
	var protein []byte
	for i := 0; i+3 <= len(upper); i += 3 {
		codon := upper[i : i+3]
		aa, ok := standardCodonTable[codon]
		if !ok {
			return "", ErrInvalidBase
		}
		protein = append(protein, aa)
	}
	return string(protein), nil
}

func TranslateDNA(dna string) (string, error) {
	if len(dna) == 0 {
		return "", ErrEmptyInput
	}
	rna := DNAToRNA(dna)
	return TranslateRNA(rna)
}

func TranslateFrame(dna string, frame int) (string, error) {
	if frame < 0 || frame > 2 {
		return "", errors.New("translate: frame must be 0, 1, or 2")
	}
	if len(dna) <= frame {
		return "", ErrEmptyInput
	}
	return TranslateDNA(dna[frame:])
}

func ThreeFrame(dna string) ([3]string, error) {
	var result [3]string
	for f := 0; f < 3; f++ {
		p, err := TranslateFrame(dna, f)
		if err != nil {
			return result, err
		}
		result[f] = p
	}
	return result, nil
}

func SixFrame(dna string) ([6]string, error) {
	var result [6]string
	fwd, err := ThreeFrame(dna)
	if err != nil {
		return result, err
	}
	result[0] = fwd[0]
	result[1] = fwd[1]
	result[2] = fwd[2]

	rc := reverseComplement(dna)
	rev, err := ThreeFrame(rc)
	if err != nil {
		return result, err
	}
	result[3] = rev[0]
	result[4] = rev[1]
	result[5] = rev[2]
	return result, nil
}

func CodonUsage(dna string) (map[string]int, error) {
	if len(dna) == 0 {
		return nil, ErrEmptyInput
	}
	upper := toUpper(dna)
	rna := DNAToRNA(upper)
	counts := make(map[string]int)
	for i := 0; i+3 <= len(rna); i += 3 {
		codon := rna[i : i+3]
		counts[codon]++
	}
	return counts, nil
}

func LookupCodon(codon string) (byte, bool) {
	aa, ok := standardCodonTable[codon]
	return aa, ok
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

func reverseComplement(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		var c byte
		switch s[len(s)-1-i] {
		case 'A', 'a':
			c = 'T'
		case 'T', 't':
			c = 'A'
		case 'C', 'c':
			c = 'G'
		case 'G', 'g':
			c = 'C'
		default:
			c = 'N'
		}
		out[i] = c
	}
	return string(out)
}

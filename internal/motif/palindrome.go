package motif

import "strings"

type Palindrome struct {
	Start  int
	End    int
	Length int
	Seq    string
}

func FindPalindromes(seq string, minLen, maxLen int) []Palindrome {
	upper := strings.ToUpper(seq)
	var results []Palindrome
	for length := minLen; length <= maxLen; length += 2 {
		for i := 0; i+length <= len(upper); i++ {
			sub := upper[i : i+length]
			if isPalindrome(sub) {
				results = append(results, Palindrome{
					Start:  i,
					End:    i + length,
					Length: length,
					Seq:    seq[i : i+length],
				})
			}
		}
	}
	return results
}

func isPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if !isComplement(s[i], s[n-1-i]) {
			return false
		}
	}
	return true
}

func isComplement(a, b byte) bool {
	switch a {
	case 'A':
		return b == 'T'
	case 'T':
		return b == 'A'
	case 'C':
		return b == 'G'
	case 'G':
		return b == 'C'
	default:
		return false
	}
}

func LongestPalindrome(palindromes []Palindrome) *Palindrome {
	if len(palindromes) == 0 {
		return nil
	}
	best := &palindromes[0]
	for i := 1; i < len(palindromes); i++ {
		if palindromes[i].Length > best.Length {
			best = &palindromes[i]
		}
	}
	return best
}

func CountByLength(palindromes []Palindrome) map[int]int {
	counts := make(map[int]int)
	for _, p := range palindromes {
		counts[p.Length]++
	}
	return counts
}

func HasHairpotential(p Palindrome) bool {
	return p.Length >= 8
}

func InvertedRepeats(seq string, armLen, minSpacer, maxSpacer int) []InvertedRepeat {
	upper := strings.ToUpper(seq)
	var results []InvertedRepeat
	for i := 0; i+armLen <= len(upper); i++ {
		arm1 := upper[i : i+armLen]
		arm1RC := revCompStr(arm1)
		for spacer := minSpacer; spacer <= maxSpacer; spacer++ {
			j := i + armLen + spacer
			if j+armLen > len(upper) {
				break
			}
			arm2 := upper[j : j+armLen]
			if arm2 == arm1RC {
				results = append(results, InvertedRepeat{
					Start1:    i,
					End1:      i + armLen,
					Start2:    j,
					End2:      j + armLen,
					ArmLength: armLen,
					Spacer:    spacer,
				})
			}
		}
	}
	return results
}

type InvertedRepeat struct {
	Start1    int
	End1      int
	Start2    int
	End2      int
	ArmLength int
	Spacer    int
}

func revCompStr(s string) string {
	out := make([]byte, len(s))
	for i := range s {
		switch s[len(s)-1-i] {
		case 'A':
			out[i] = 'T'
		case 'T':
			out[i] = 'A'
		case 'C':
			out[i] = 'G'
		case 'G':
			out[i] = 'C'
		default:
			out[i] = 'N'
		}
	}
	return string(out)
}

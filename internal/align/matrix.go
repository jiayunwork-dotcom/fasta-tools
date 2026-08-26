package align

type SubMatrix struct {
	scores map[byte]map[byte]int
	name   string
}

func NewSubMatrix(name string, scores map[byte]map[byte]int) *SubMatrix {
	return &SubMatrix{name: name, scores: scores}
}

func (m *SubMatrix) Score(a, b byte) int {
	if row, ok := m.scores[a]; ok {
		if s, ok := row[b]; ok {
			return s
		}
	}
	if a == b {
		return 1
	}
	return -1
}

func (m *SubMatrix) Name() string { return m.name }

func DNAFullMatrix() *SubMatrix {
	bases := []byte{'A', 'C', 'G', 'T'}
	scores := make(map[byte]map[byte]int)
	for _, a := range bases {
		scores[a] = make(map[byte]int)
		for _, b := range bases {
			if a == b {
				scores[a][b] = 5
			} else if isTransition(a, b) {
				scores[a][b] = -1
			} else {
				scores[a][b] = -4
			}
		}
	}
	return &SubMatrix{name: "DNA-Full", scores: scores}
}

func isTransition(a, b byte) bool {
	purines := map[byte]bool{'A': true, 'G': true}
	pyrimidines := map[byte]bool{'C': true, 'T': true}
	return (purines[a] && purines[b]) || (pyrimidines[a] && pyrimidines[b])
}

func IdentityMatrix() *SubMatrix {
	bases := []byte{'A', 'C', 'G', 'T', 'N'}
	scores := make(map[byte]map[byte]int)
	for _, a := range bases {
		scores[a] = make(map[byte]int)
		for _, b := range bases {
			if a == b {
				scores[a][b] = 1
			}
		}
	}
	return &SubMatrix{name: "Identity", scores: scores}
}

func GlobalWithMatrix(a, b string, m *SubMatrix, gapOpen int) (*Result, error) {
	if len(a) == 0 || len(b) == 0 {
		return nil, ErrEmptySequence
	}
	la, lb := len(a), len(b)
	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
	}
	for i := 1; i <= la; i++ {
		dp[i][0] = dp[i-1][0] + gapOpen
	}
	for j := 1; j <= lb; j++ {
		dp[0][j] = dp[0][j-1] + gapOpen
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			s := m.Score(a[i-1], b[j-1])
			diag := dp[i-1][j-1] + s
			up := dp[i-1][j] + gapOpen
			left := dp[i][j-1] + gapOpen
			v, _ := max3(diag, up, left)
			dp[i][j] = v
		}
	}
	var alignA, alignB []byte
	i, j := la, lb
	for i > 0 || j > 0 {
		if i > 0 && j > 0 {
			s := m.Score(a[i-1], b[j-1])
			if dp[i][j] == dp[i-1][j-1]+s {
				alignA = append(alignA, a[i-1])
				alignB = append(alignB, b[j-1])
				i--
				j--
				continue
			}
		}
		if i > 0 && dp[i][j] == dp[i-1][j]+gapOpen {
			alignA = append(alignA, a[i-1])
			alignB = append(alignB, '-')
			i--
		} else {
			alignA = append(alignA, '-')
			alignB = append(alignB, b[j-1])
			j--
		}
	}
	reverse(alignA)
	reverse(alignB)
	res := &Result{
		AlignedA: string(alignA),
		AlignedB: string(alignB),
		Score:    dp[la][lb],
		Length:   len(alignA),
	}
	res.Identity = computeIdentity(alignA, alignB)
	return res, nil
}

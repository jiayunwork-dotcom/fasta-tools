package translate

type GeneticCode struct {
	Name  string
	ID    int
	Table map[string]byte
}

func MitochondrialCode() *GeneticCode {
	table := make(map[string]byte)
	for k, v := range standardCodonTable {
		table[k] = v
	}
	table["UGA"] = 'W'
	table["AGA"] = '*'
	table["AGG"] = '*'
	table["AUA"] = 'M'
	return &GeneticCode{Name: "Vertebrate Mitochondrial", ID: 2, Table: table}
}

func YeastMitoCode() *GeneticCode {
	table := make(map[string]byte)
	for k, v := range standardCodonTable {
		table[k] = v
	}
	table["CUU"] = 'T'
	table["CUC"] = 'T'
	table["CUA"] = 'T'
	table["CUG"] = 'T'
	table["UGA"] = 'W'
	return &GeneticCode{Name: "Yeast Mitochondrial", ID: 3, Table: table}
}

func TranslateWithCode(rna string, code *GeneticCode) (string, error) {
	if len(rna) == 0 {
		return "", ErrEmptyInput
	}
	upper := toUpper(rna)
	var protein []byte
	for i := 0; i+3 <= len(upper); i += 3 {
		codon := upper[i : i+3]
		aa, ok := code.Table[codon]
		if !ok {
			return "", ErrInvalidBase
		}
		protein = append(protein, aa)
	}
	return string(protein), nil
}

func ReverseTranslate(aminoAcid byte) []string {
	var codons []string
	for codon, aa := range standardCodonTable {
		if aa == aminoAcid {
			codons = append(codons, codon)
		}
	}
	return codons
}

func IsStartCodon(codon string) bool {
	return toUpper(codon) == "ATG"
}

func IsStopCodon(codon string) bool {
	upper := toUpper(codon)
	rna := DNAToRNA(upper)
	return rna == "UAA" || rna == "UAG" || rna == "UGA"
}

func CodonCount(seq string) int {
	return len(seq) / 3
}

func Degenerate(aminoAcid byte) int {
	return len(ReverseTranslate(aminoAcid))
}

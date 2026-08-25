package digest

import (
	"errors"
	"sort"
	"strings"
)

var ErrNoEnzyme = errors.New("digest: no enzyme specified")

var ErrEmptySeq = errors.New("digest: empty sequence")

type Enzyme struct {
	Name   string
	Site   string
	CutPos int
}

type Fragment struct {
	Start  int
	End    int
	Length int
	Seq    string
}

func CommonEnzymes() map[string]Enzyme {
	return map[string]Enzyme{
		"EcoRI":   {Name: "EcoRI", Site: "GAATTC", CutPos: 1},
		"BamHI":   {Name: "BamHI", Site: "GGATCC", CutPos: 1},
		"HindIII": {Name: "HindIII", Site: "AAGCTT", CutPos: 1},
		"NotI":    {Name: "NotI", Site: "GCGGCCGC", CutPos: 2},
		"XhoI":    {Name: "XhoI", Site: "CTCGAG", CutPos: 1},
		"SalI":    {Name: "SalI", Site: "GTCGAC", CutPos: 1},
		"PstI":    {Name: "PstI", Site: "CTGCAG", CutPos: 5},
		"SmaI":    {Name: "SmaI", Site: "CCCGGG", CutPos: 3},
		"KpnI":    {Name: "KpnI", Site: "GGTACC", CutPos: 5},
		"SacI":    {Name: "SacI", Site: "GAGCTC", CutPos: 5},
	}
}

type CutSite struct {
	Enzyme    string
	Position  int
	SiteStart int
}

func FindSites(seq string, enz Enzyme) ([]CutSite, error) {
	if len(seq) == 0 {
		return nil, ErrEmptySeq
	}
	upper := strings.ToUpper(seq)
	site := strings.ToUpper(enz.Site)
	var sites []CutSite
	for i := 0; i+len(site) <= len(upper); i++ {
		if upper[i:i+len(site)] == site {
			sites = append(sites, CutSite{
				Enzyme:    enz.Name,
				Position:  i + enz.CutPos,
				SiteStart: i,
			})
		}
	}
	return sites, nil
}

func Digest(seq string, enz Enzyme) ([]Fragment, error) {
	if len(seq) == 0 {
		return nil, ErrEmptySeq
	}
	sites, err := FindSites(seq, enz)
	if err != nil {
		return nil, err
	}
	if len(sites) == 0 {
		return []Fragment{{Start: 0, End: len(seq), Length: len(seq), Seq: seq}}, nil
	}
	cuts := make([]int, 0, len(sites))
	for _, s := range sites {
		cuts = append(cuts, s.Position)
	}
	sort.Ints(cuts)
	var frags []Fragment
	prev := 0
	for _, c := range cuts {
		if c > prev && c <= len(seq) {
			frags = append(frags, Fragment{
				Start:  prev,
				End:    c,
				Length: c - prev,
				Seq:    seq[prev:c],
			})
			prev = c
		}
	}
	if prev < len(seq) {
		frags = append(frags, Fragment{
			Start:  prev,
			End:    len(seq),
			Length: len(seq) - prev,
			Seq:    seq[prev:],
		})
	}
	return frags, nil
}

func DigestMulti(seq string, enzymes []Enzyme) ([]Fragment, error) {
	if len(seq) == 0 {
		return nil, ErrEmptySeq
	}
	if len(enzymes) == 0 {
		return nil, ErrNoEnzyme
	}
	var allCuts []int
	for _, enz := range enzymes {
		sites, err := FindSites(seq, enz)
		if err != nil {
			return nil, err
		}
		for _, s := range sites {
			allCuts = append(allCuts, s.Position)
		}
	}
	sort.Ints(allCuts)
	if len(allCuts) > 0 {
		j := 0
		for i := 1; i < len(allCuts); i++ {
			if allCuts[i] != allCuts[j] {
				j++
				allCuts[j] = allCuts[i]
			}
		}
		allCuts = allCuts[:j+1]
	}
	var frags []Fragment
	prev := 0
	for _, c := range allCuts {
		if c > prev && c <= len(seq) {
			frags = append(frags, Fragment{
				Start:  prev,
				End:    c,
				Length: c - prev,
				Seq:    seq[prev:c],
			})
			prev = c
		}
	}
	if prev < len(seq) {
		frags = append(frags, Fragment{
			Start:  prev,
			End:    len(seq),
			Length: len(seq) - prev,
			Seq:    seq[prev:],
		})
	}
	return frags, nil
}

func FragmentLengths(frags []Fragment) []int {
	lengths := make([]int, len(frags))
	for i, f := range frags {
		lengths[i] = f.Length
	}
	sort.Sort(sort.Reverse(sort.IntSlice(lengths)))
	return lengths
}

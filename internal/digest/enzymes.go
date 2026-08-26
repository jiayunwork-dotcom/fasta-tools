package digest

import "strings"

type EnzymeDB struct {
	entries map[string]Enzyme
}

func NewEnzymeDB() *EnzymeDB {
	db := &EnzymeDB{entries: make(map[string]Enzyme)}
	for name, enz := range CommonEnzymes() {
		db.entries[name] = enz
	}
	return db
}

func (db *EnzymeDB) Add(e Enzyme) {
	db.entries[e.Name] = e
}

func (db *EnzymeDB) Get(name string) (Enzyme, bool) {
	for n, e := range db.entries {
		if strings.EqualFold(n, name) {
			return e, true
		}
	}
	return Enzyme{}, false
}

func (db *EnzymeDB) SearchBySite(sub string) []Enzyme {
	upper := strings.ToUpper(sub)
	var matches []Enzyme
	for _, e := range db.entries {
		if strings.Contains(e.Site, upper) {
			matches = append(matches, e)
		}
	}
	return matches
}

func (db *EnzymeDB) ListAll() []Enzyme {
	out := make([]Enzyme, 0, len(db.entries))
	for _, e := range db.entries {
		out = append(out, e)
	}
	return out
}

func (db *EnzymeDB) Count() int { return len(db.entries) }

func (db *EnzymeDB) Remove(name string) {
	delete(db.entries, name)
}

func (db *EnzymeDB) FindCutters(seq string) []Enzyme {
	var cutters []Enzyme
	for _, e := range db.entries {
		sites, err := FindSites(seq, e)
		if err == nil && len(sites) > 0 {
			cutters = append(cutters, e)
		}
	}
	return cutters
}

func (db *EnzymeDB) FindNonCutters(seq string) []Enzyme {
	var nonCutters []Enzyme
	for _, e := range db.entries {
		sites, err := FindSites(seq, e)
		if err == nil && len(sites) == 0 {
			nonCutters = append(nonCutters, e)
		}
	}
	return nonCutters
}

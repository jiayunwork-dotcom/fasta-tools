package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"fasta-tools/internal/fasta"
	"fasta-tools/internal/kmer"
	"fasta-tools/internal/seq"
)

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	addr := fs.String("addr", ":8080", "listen address")
	if err := fs.Parse(reorder(args)); err != nil {
		os.Exit(2)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/gc", handleGC)
	mux.HandleFunc("/api/rc", handleRC)
	mux.HandleFunc("/api/kmer", handleKmer)
	mux.HandleFunc("/api/parse", handleParse)

	fmt.Printf("fasta-tools serving on %s\n", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		fmt.Fprintf(os.Stderr, "serve: %v\n", err)
		os.Exit(1)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleGC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	records, err := fasta.Parse(r.Body)
	if err != nil {
		http.Error(w, "parse fasta: "+err.Error(), http.StatusBadRequest)
		return
	}
	type result struct {
		Header string  `json:"header"`
		GC     float64 `json:"gc_content"`
		Length int     `json:"length"`
	}
	var results []result
	for _, rec := range records {
		results = append(results, result{
			Header: rec.Header,
			GC:     seq.GCContent(rec.Sequence),
			Length: len(rec.Sequence),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleRC(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	records, err := fasta.Parse(r.Body)
	if err != nil {
		http.Error(w, "parse fasta: "+err.Error(), http.StatusBadRequest)
		return
	}
	type result struct {
		Header string `json:"header"`
		RC     string `json:"reverse_complement"`
	}
	var results []result
	for _, rec := range records {
		results = append(results, result{
			Header: rec.Header,
			RC:     seq.ReverseComplement(rec.Sequence),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func handleKmer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	kStr := r.URL.Query().Get("k")
	k := 3
	if kStr != "" {
		fmt.Sscanf(kStr, "%d", &k)
	}
	records, err := fasta.Parse(r.Body)
	if err != nil {
		http.Error(w, "parse fasta: "+err.Error(), http.StatusBadRequest)
		return
	}
	var seqs []string
	for _, rec := range records {
		seqs = append(seqs, rec.Sequence)
	}
	combined := strings.Join(seqs, "")
	counts, err := kmer.Count(combined, k)
	if err != nil {
		http.Error(w, "kmer count: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(counts)
}

func handleParse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	records, err := fasta.Parse(r.Body)
	if err != nil {
		http.Error(w, "parse fasta: "+err.Error(), http.StatusBadRequest)
		return
	}
	type info struct {
		Header string `json:"header"`
		Length int    `json:"length"`
	}
	var out []info
	for _, rec := range records {
		out = append(out, info{Header: rec.Header, Length: len(rec.Sequence)})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

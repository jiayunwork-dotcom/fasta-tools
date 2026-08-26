package fasta

import (
	"fmt"
	"io"
	"strings"
)

const DefaultLineWidth = 80

type Writer struct {
	w         io.Writer
	lineWidth int
	written   int
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w, lineWidth: DefaultLineWidth}
}

func NewWriterWidth(w io.Writer, width int) *Writer {
	if width <= 0 {
		width = DefaultLineWidth
	}
	return &Writer{w: w, lineWidth: width}
}

func (fw *Writer) Write(r Record) error {
	if _, err := fmt.Fprintf(fw.w, ">%s\n", r.Header); err != nil {
		return err
	}
	seq := r.Sequence
	for i := 0; i < len(seq); i += fw.lineWidth {
		end := i + fw.lineWidth
		if end > len(seq) {
			end = len(seq)
		}
		if _, err := fmt.Fprintf(fw.w, "%s\n", seq[i:end]); err != nil {
			return err
		}
	}
	if len(seq) == 0 {
		if _, err := fmt.Fprint(fw.w, "\n"); err != nil {
			return err
		}
	}
	fw.written++
	return nil
}

func (fw *Writer) WriteAll(records []Record) error {
	for _, r := range records {
		if err := fw.Write(r); err != nil {
			return err
		}
	}
	return nil
}

func (fw *Writer) Written() int { return fw.written }

func (fw *Writer) LineWidth() int { return fw.lineWidth }

func Merge(records []Record, header string) Record {
	var sb strings.Builder
	for _, r := range records {
		sb.WriteString(r.Sequence)
	}
	return Record{Header: header, Sequence: sb.String()}
}

func Split(r Record, chunkSize int) []Record {
	if chunkSize <= 0 {
		chunkSize = 1
	}
	var out []Record
	for i := 0; i < len(r.Sequence); i += chunkSize {
		end := i + chunkSize
		if end > len(r.Sequence) {
			end = len(r.Sequence)
		}
		hdr := fmt.Sprintf("%s_%d", r.Header, len(out)+1)
		out = append(out, Record{Header: hdr, Sequence: r.Sequence[i:end]})
	}
	return out
}

func Subset(records []Record, indices []int) []Record {
	var out []Record
	for _, idx := range indices {
		if idx >= 0 && idx < len(records) {
			out = append(out, records[idx])
		}
	}
	return out
}

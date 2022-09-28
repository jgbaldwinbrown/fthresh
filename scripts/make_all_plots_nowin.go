package main

import (
	"os"
	"io"
	"bufio"
	"github.com/jgbaldwinbrown/makem"
)

func AddEntry(m *makem.MakeData, s string) {
	main := s + ".txt"
	bed := s + ".bed"
	fdr := s + "_fdr.bed"
	plfmt := s + "_fdr_plfmt.bed"

	r := makem.Recipe{}
	r.AddTargets(bed)
	r.AddDeps(main)
	r.AddScripts("bedify " + main + " > " + bed)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(fdr)
	r.AddDeps(bed)
	r.AddScripts("fdr_it <" + bed + " > " + fdr)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(plfmt)
	r.AddDeps(fdr)
	r.AddScripts("plfmt <" + fdr + " > " + plfmt)
	m.Add(r)
}

func MakeMakefile(r io.Reader, w io.Writer) {
	makefile := new(makem.MakeData)

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)
	for s.Scan() {
		AddEntry(makefile, s.Text())
	}

	makefile.Fprint(w)
}

func main() {
	MakeMakefile(os.Stdin, os.Stdout)
}

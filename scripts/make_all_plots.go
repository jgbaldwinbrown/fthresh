package main

import (
	"fmt"
	"flag"
	"os"
	"io"
	"bufio"
	"github.com/jgbaldwinbrown/makem"
)

func AddEntry(m *makem.MakeData, s string, f Flags) {
	main := s + ".txt"
	win := s + "_win.txt"
	winbed := s + "_win.bed"
	winfdr := s + "_win_fdr.bed"
	plfmt := s + "_win_fdr_plfmt.bed"
	plot := s + "_win_fdr_plot.png"

	if f.WinSize != 50000 || f.WinStep != 5000 {
		winprefix := fmt.Sprintf("%s_win%d_%d", s, f.WinSize, f.WinStep)
		win = winprefix + ".txt"
		winbed = winprefix + ".bed"
		winfdr = winprefix + "_fdr.bed"
		plfmt = winprefix + "_fdr_plfmt.bed"
		plot = winprefix + "_fdr_plot.png"
	}

	r := makem.Recipe{}
	r.AddTargets(win)
	r.AddDeps(main)
	winscript := fmt.Sprintf("python3 window_fisher_bp.py <%s 2 0 1 %d %d > %s", main, f.WinSize, f.WinStep, win)
	r.AddScripts(winscript)
	// r.AddScripts("python3 window_fisher_bp.py <" + main + " 2 0 1 50000 5000 > " + win)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(winbed)
	r.AddDeps(win)
	r.AddScripts("bedify " + win + " > " + winbed)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(winfdr)
	r.AddDeps(winbed)
	r.AddScripts("fdr_it <" + winbed + " > " + winfdr)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(plfmt)
	r.AddDeps(winfdr)
	r.AddScripts("plfmt <" + winfdr + " > " + plfmt)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(plot)
	r.AddDeps(plfmt)
	r.AddScripts("Rscript plot_pretty_hlines_bp.R " + plfmt + " " + plot)
	m.Add(r)
}

func MakeMakefile(r io.Reader, w io.Writer, flags Flags) {
	makefile := new(makem.MakeData)

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)
	for s.Scan() {
		AddEntry(makefile, s.Text(), flags)
	}

	makefile.Fprint(w)
}

type Flags struct {
	WinSize int
	WinStep int
}

func GetFlags() Flags {
	var f Flags
	flag.IntVar(&f.WinSize, "w", 50000, "Window size")
	flag.IntVar(&f.WinStep, "s", 5000, "Window size")
	flag.Parse()
	return f
}

func main() {
	flags := GetFlags()
	MakeMakefile(os.Stdin, os.Stdout, flags)
}

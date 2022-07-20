package main

import (
	"os"
	"io"
	"bufio"
	"github.com/jgbaldwinbrown/makem"
)

func AddEntry(m *makem.MakeData, s string) {
	main := s + ".txt"
	win := s + "_win.txt"
	winbed := s + "_win.bed"
	winfdr := s + "_win_fdr.bed"
	plfmt := s + "_win_fdr_plfmt.bed"
	plot := s + "_win_fdr_plot.png"

	r := makem.Recipe{}
	r.AddTargets(win)
	r.AddDeps(main)
	r.AddScripts("python3 window_fisher_bp.py <" + main + " 2 0 1 50000 5000 > " + win)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(winbed)
	r.AddDeps(win)
	r.AddScripts("bash bedify.sh " + win + " > " + winbed)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(winfdr)
	r.AddDeps(winbed)
	r.AddScripts("./fdr_it.py <" + winbed + " > " + winfdr)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(plfmt)
	r.AddDeps(winfdr)
	r.AddScripts("./plfmt.py <" + winfdr + " > " + plfmt)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(plot)
	r.AddDeps(plfmt)
	r.AddScripts("Rscript plot_pretty_hlines_bp.R " + plfmt + " " + plot)
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

/*
all: color_af_thresh_pfst_win.txt color_af_thresh_pfst_win.bed color_af_thresh_pfst_win_fdr.bed color_fdr_plfmt.bed color_pretty_plot_bp.png color_pretty_plot_bp_hlines.png

SHELL := /bin/bash

.PHONY: all

.DELETE_ON_ERROR:

color_af_thresh_pfst_win.txt: color_af_thresh_pfst.txt
	python3 window_fisher.py <color_af_thresh_pfst.txt 2 10 3 > color_af_thresh_pfst_win.txt

color_af_thresh_pfst_win.bed: color_af_thresh_pfst_win.txt
	bash bedify.sh color_af_thresh_pfst_win.txt > color_af_thresh_pfst_win.bed

color_af_thresh_pfst_win_fdr.bed: color_af_thresh_pfst_win.bed
	./fdr_it.py <color_af_thresh_pfst_win.bed > color_af_thresh_pfst_win_fdr.bed

color_fdr_plfmt.bed: color_af_thresh_pfst_win_fdr.bed
	./plfmt.py <$^ > $@

color_pretty_plot_bp.png: plot_pretty_bp.R color_fdr_plfmt.bed
	Rscript $^ $@

color_pretty_plot_bp_hlines.png: plot_pretty_hlines_bp.R color_fdr_plfmt.bed
	Rscript $^ $@
*/

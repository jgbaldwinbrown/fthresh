package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

func Preamble() string {
	return `\documentclass{article}

%%% Load packages
\usepackage[margin=1in]{geometry}
\usepackage{amsthm,amsmath}
\RequirePackage{natbib}
\RequirePackage{hyperref}
\usepackage[T1]{fontenc}
\usepackage[utf8]{inputenc} %unicode support
\usepackage[polutonikogreek,english]{babel}
\usepackage{listings}
\usepackage{graphicx}

\usepackage{longtable}

\begin{document}
`
}

func FigStart() string {
	return `\begin{figure}[h!]`
}

func WrapGraphics(in string) string {
	return `\includegraphics[width=\linewidth]{` + in + `}`
}

func WrapCaption(bold, normal string) string {
	return strings.Replace(`\caption{\textbf{` + bold + `} ` + normal + `}`, "_", "\\_", -1)
}

func CleanHomer(homer string) string {
	return strings.Replace(strings.ToLower(homer), "homer", " homer", -1)
}

func GenCaption(path string) string {
	split := strings.Split(path, "_")
	b1 := CleanHomer(split[2])
	t1 := split[4]
	bit1 := strings.ToLower(split[6])
	r1 := split[8]

	b2 := CleanHomer(split[10])
	t2 := split[12]
	bit2 := strings.ToLower(split[14])
	r2 := split[16]

	trait := strings.ToLower(split[17])
	cut1 := strings.ToLower(split[18])
	cut2 := strings.ToLower(split[19])

	return fmt.Sprint(
		"A comparison of pFst, Fst, and selection coefficient when comparing replicate " +
		r1 +
		" of " +
		cut1 +
		" " +
		bit1 +
		" " +
		b1 +
		" lice at time point " +
		t1 +
		" to replicate " +
		r2 +
		" of " +
		cut2 +
		" " +
		bit2 +
		" " +
		b2 +
		" lice at time point " +
		t2 +
		" in terms of " +
		trait + ".")
}
// _breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Low_Mid___multiplot_plot_pfst_fst_selec_v2.png

func FigEnd() string {
	return `\end{figure}`
}

func Ending() string {
	return `\end{document}`
}

func main() {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	fmt.Fprintln(w, Preamble())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Fprintln(w, FigStart())
		fmt.Fprintln(w, WrapGraphics(s.Text()))
		fmt.Fprintln(w, WrapCaption("", GenCaption(s.Text())))
		fmt.Fprintln(w, FigEnd())
	}
	fmt.Fprintln(w, Ending())
}

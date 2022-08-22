package fthresh

// _pfst_plfmt.bed

import (
	"fmt"
	"os"
	"io"
)

func FprintPlotTogether(w io.Writer, pfstpaths, sigpaths []string) {
	fmt.Fprintf(w, `#!/bin/bash
set -e

plot_goods_pfst_together \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
> together_out.txt
`,
		pfstpaths[0],
		sigpaths[0],
		pfstpaths[1],
		sigpaths[1],
		pfstpaths[2],
		sigpaths[2],
		pfstpaths[3],
		sigpaths[3],
		"good4plots_together.png",
	)
}

func PlotSetToTogetherScript() {
	plotsets := ReadPlotSets(os.Stdin)
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, set.GoodPfstSpans)
	}
	FprintPlotTogether(os.Stdout, pfstpaths, sigpaths)
}

func PlotSetToTogetherSubFullScript() {
	plotsets := ReadPlotSets(os.Stdin)
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, SubFullPath(set.GoodPfstSpans))
	}
	FprintPlotTogether(os.Stdout, pfstpaths, sigpaths)
}

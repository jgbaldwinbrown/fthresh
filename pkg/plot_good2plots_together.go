package fthresh

// _pfst_plfmt.bed

import (
	"fmt"
	"os"
	"io"
)

func FprintPlotTogether2(w io.Writer, pfstpaths, sigpaths []string, out string) {
	fmt.Printf(`#!/bin/bash
set -e

plot_goods_pfst_together2 \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
> %v_out.txt
`,
		pfstpaths[0],
		sigpaths[0],
		pfstpaths[1],
		sigpaths[1],
		out,
		out,
	)
}

func PlotSetToTogetherScript2() {
	plotsets := ReadPlotSets(os.Stdin)
	PlotSetToTogetherSubSimpleScript2(plotsets)
	PlotSetToTogetherSubFullScript2(plotsets)
}

func PlotSetToTogetherSubSimpleScript2(plotsets []PlotSet) {
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, set.GoodPfstSpans)
	}
	FprintPlotTogether2(os.Stdout, pfstpaths, sigpaths, "good2plots_together.png")
}

func PlotSetToTogetherSubFullScript2(plotsets []PlotSet) {
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, SubFullPath(set.GoodPfstSpans))
	}
	FprintPlotTogether2(os.Stdout, pfstpaths, sigpaths, "good2plots_together_subfull.png")
}

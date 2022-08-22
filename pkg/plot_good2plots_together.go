package fthresh

// _pfst_plfmt.bed

import (
	"fmt"
	"os"
)

func PlotSetToTogetherScript2() {
	plotsets := ReadPlotSets(os.Stdin)
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, set.GoodPfstSpans)
	}
	fmt.Printf(`#!/bin/bash
set -e

plot_goods_pfst_together2 \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
> together2_out.txt
`,
		pfstpaths[0],
		sigpaths[0],
		pfstpaths[1],
		sigpaths[1],
		"good2plots_together.png",
	)
}

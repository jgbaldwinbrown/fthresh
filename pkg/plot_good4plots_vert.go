package fthresh

// _pfst_plfmt.bed

import (
	"fmt"
	"os"
)

func PlotSetToVertScript() {
	plotsets := ReadPlotSets(os.Stdin)
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, set.GoodPfstSpans)
	}
	fmt.Printf(`#!/bin/bash
set -e

plot_goods_pfst_vert \
	%v \
	%v \
	%v \
	%v \
	%v \
	%v \
	%v \
	%v \
	%v \
> vert_out.txt
`,
		pfstpaths[0],
		sigpaths[0],
		pfstpaths[1],
		sigpaths[1],
		pfstpaths[2],
		sigpaths[2],
		pfstpaths[3],
		sigpaths[3],
		"good4plots_vert.png",
	)
}

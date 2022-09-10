package fthresh

// _pfst_plfmt.bed

import (
	"fmt"
	"os"
)

func PlotSetToVertScript8() {
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
> vert_out_color.txt

plot_goods_pfst_vert_size \
	%v \
	%v \
	%v \
	%v \
	%v \
	%v \
	%v \
	%v \
	%v \
> vert_out_size.txt
`,
		pfstpaths[0],
		sigpaths[0],
		pfstpaths[1],
		sigpaths[1],
		pfstpaths[2],
		sigpaths[2],
		pfstpaths[3],
		sigpaths[3],
		"good8plots_vert_color.png",
		pfstpaths[4],
		sigpaths[4],
		pfstpaths[5],
		sigpaths[5],
		pfstpaths[6],
		sigpaths[6],
		pfstpaths[7],
		sigpaths[7],
		"good8plots_vert_size.png",
	)
}

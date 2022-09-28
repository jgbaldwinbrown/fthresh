package fthresh

import (
	"fmt"
	// "sync"
	"runtime"
	"os"
)

func PrebuiltGoodAndAltsReps() []GoodAndAlts {
	var ga []GoodAndAlts
	for _, breed := range []string{"black", "white", "figurita", "runt"} {
		for _, rep := range []int{1,2,3,4} {
			ga = append(ga, GoodAndAlts{
					Good: Comp{breed, "unbitted", rep, 36, "feral", "unbitted", rep, 36},
					Alts: []Comp {
						Comp{breed, "bitted", rep, 36, "feral", "bitted", rep, 36},
					},
				},
			)
		}
	}
	return ga
}

func RunGood4PlotsReps() {
	flags := GetPlot4Flags()
	plot_sets := ReadPlotSets(os.Stdin)
	runtime.GOMAXPROCS(flags.Threads)

	errors := PercTMASets(plot_sets, flags.Percentile)
	for _, err := range errors {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}

	goodsAndAlts := PrebuiltGoodAndAltsReps()
	statistics := []string{"pFst", "Fst", "Selec"}
	outpaths, errors := SubtractAllAlts(goodsAndAlts, statistics, false, flags.FullReps, flags.Threads)
	for _, err := range errors {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	pathsconn, err := os.Create("goodbedpaths_rep.txt")
	if err != nil { panic(err) }
	defer pathsconn.Close()
	for _, path := range outpaths {
		fmt.Fprintln(pathsconn, path)
	}

	outpaths_f, errors_f := SubtractAllAlts(goodsAndAlts, statistics, true, flags.FullReps, flags.Threads)
	for _, err := range errors_f {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}

	pathsconn_f, err_f := os.Create("goodbedpaths_rep_f.txt")
	if err_f != nil { panic(err_f) }
	defer pathsconn_f.Close()
	for _, path := range outpaths_f {
		fmt.Fprintln(pathsconn_f, path)
	}
}

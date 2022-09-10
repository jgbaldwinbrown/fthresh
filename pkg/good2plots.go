package fthresh

import (
	"log"
	"fmt"
	"runtime"
	"os"
)

func PrebuiltGoodAndAlts2() []GoodAndAlts {
	return []GoodAndAlts {
		GoodAndAlts{
			Good: Comp{"black", "unbitted", 0, 36, "feral", "unbitted", 0, 36},
			Alts: []Comp {
				Comp{"black", "bitted", 0, 36, "feral", "bitted", 0, 36},
			},
		},
		GoodAndAlts{
			Good: Comp{"white", "unbitted", 0, 36, "feral", "unbitted", 0, 36},
			Alts: []Comp {
				Comp{"white", "bitted", 0, 36, "feral", "bitted", 0, 36},
			},
		},
	}
}

func RunGood2Plots() {
	flags := GetPlot4Flags()
	plot_sets := ReadPlotSets(os.Stdin)
	runtime.GOMAXPROCS(flags.Threads)

	fmt.Fprintf(os.Stderr, "g one\n")
	errors := PercTMASets(plot_sets, flags.Percentile)
	for _, err := range errors {
		if err != nil { log.Print(err) }
	}
	fmt.Fprintf(os.Stderr, "g two\n")

	goodsAndAlts := PrebuiltGoodAndAlts2()
	fmt.Println(goodsAndAlts)
	// statistics := []string{"pFst", "Fst", "Selec"}
	statistics := []string{"pFst", "Fst"}
	outpaths, errors := SubtractAllAlts(goodsAndAlts, statistics, false, flags.Threads)
	for _, err := range errors {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}
	fmt.Fprintf(os.Stderr, "g three\n")

	outpaths_f, errors_f := SubtractAllAlts(goodsAndAlts, statistics, true, flags.Threads)
	for _, err := range errors_f {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}
	fmt.Fprintf(os.Stderr, "g four\n")

	pathsconn, err := os.Create("goodbedpaths2.txt")
	if err != nil { panic(err) }
	defer pathsconn.Close()
	for _, path := range outpaths {
		fmt.Fprintln(pathsconn, path)
	}
	fmt.Fprintf(os.Stderr, "g five\n")

	pathsconn_f, err_f := os.Create("goodbedpaths2_f.txt")
	if err_f != nil { panic(err_f) }
	defer pathsconn_f.Close()
	for _, path := range outpaths_f {
		fmt.Fprintln(pathsconn_f, path)
	}
	fmt.Fprintf(os.Stderr, "g six\n")

}

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
			Good: Comp{"black", "unbitted", "feral", "unbitted"},
			Alts: []Comp {
				Comp{"black", "bitted", "feral", "bitted"},
			},
		},
		GoodAndAlts{
			Good: Comp{"white", "unbitted", "feral", "unbitted"},
			Alts: []Comp {
				Comp{"white", "bitted", "feral", "bitted"},
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
	outpaths, errors := SubtractAllAlts(goodsAndAlts, statistics, flags.Threads)
	for _, err := range errors {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}
	fmt.Fprintf(os.Stderr, "g three\n")

	pathsconn, err := os.Create("goodbedpaths2.txt")
	if err != nil { panic(err) }
	defer pathsconn.Close()
	for _, path := range outpaths {
		fmt.Fprintln(pathsconn, path)
	}

}

package fthresh

// _pfst_plfmt.bed

import (
	"sort"
	"flag"
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

func FprintPlotTogetherNowin(w io.Writer, pfstpaths, sigpaths []string) {
	fmt.Fprintf(w, `#!/bin/bash
set -e

plot_goods_pfst_together_nowin \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
> together_out_nowin.txt
`,
		pfstpaths[0],
		sigpaths[0],
		pfstpaths[1],
		sigpaths[1],
		pfstpaths[2],
		sigpaths[2],
		pfstpaths[3],
		sigpaths[3],
		"good4plots_together_nowin.png",
	)
}

func PlotSetToTogetherScript() {
	plotsets := ReadAllrepExpPlotSets(os.Stdin)
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, set.GoodPfstSpans)
	}
	FprintPlotTogether(os.Stdout, pfstpaths, sigpaths)
}

func PlotSetToTogetherSubFullScript() {
	plotsets := ReadAllrepExpPlotSets(os.Stdin)
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, set.GoodPfstSpans)
	}
	FprintPlotTogether(os.Stdout, pfstpaths, sigpaths)
}

func psIndex(rep string, breed string, bit string) int {
	m := map[string]int {
		"R1": 0,
		"R2": 1,
		"R3": 2,
		"R4": 3,
		"black": 0,
		"white": 1,
		"figurita": 2,
		"runt": 3,
		"unbitted": 0,
		"bitted": 1,
	}
	return 8 * m[rep] + 2 * m[breed] + m[bit]
}

func psIndexTop(rep string, breed string, bit string) int {
	m := map[string]int {
		"All": 0,
		"R1": 1,
		"R2": 2,
		"R3": 3,
		"R4": 4,
		"black": 0,
		"white": 1,
		"figurita": 2,
		"runt": 3,
		"unbitted": 0,
		"bitted": 1,
	}
	// return 8 * m[rep] + 2 * m[breed] + m[bit]
	return 4 * m[rep] + m[breed]
}

func psIndexReps(breed string, bit string) []int {
	out := make([]int, 4)
	for i, rep := range []string{"R1", "R2", "R3", "R4"} {
		out[i] = psIndex(rep, breed, bit)
	}
	return out
}

func psIndexRepsTop(breed string, bit string) []int {
	reps := []string{"All", "R1", "R2", "R3", "R4"}
	out := make([]int, len(reps))
	for i, rep := range reps {
		out[i] = psIndexTop(rep, breed, bit)
	}
	return out
}

func FprintPlotTogetherRepsBreed(w io.Writer, pfstpaths, sigpaths []string, breed string) {
	idxs := psIndexReps(breed, "unbitted")
	fmt.Fprintf(w, `plot_goods_pfst_reps_together \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	good4plots_together_reps_%v.png \
> together_reps_out_%v.txt
`,
		pfstpaths[idxs[0]],
		sigpaths[idxs[0]],
		pfstpaths[idxs[1]],
		sigpaths[idxs[1]],
		pfstpaths[idxs[2]],
		sigpaths[idxs[2]],
		pfstpaths[idxs[3]],
		sigpaths[idxs[3]],
		breed,
		breed,
	)
}

func FprintPlotTogetherReps(w io.Writer, pfstpaths, sigpaths []string) {
	fmt.Fprintf(w, `#!/bin/bash
`)
	for _, breed := range []string{"black", "white", "figurita", "runt"} {
		fmt.Fprintf(w, "\n")
		FprintPlotTogetherRepsBreed(w, pfstpaths, sigpaths, breed)
	}
}

func PlotSetToTogetherSubFullScriptReps() {
	plotsets := ReadEveryRepExpPlotSets(os.Stdin)
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, set.GoodPfstSpans)
	}
	FprintPlotTogetherReps(os.Stdout, pfstpaths, sigpaths)
}

func FprintPlotTogetherRepsTopBreed(w io.Writer, pfstpaths, sigpaths []string, breed string, fullreps, nowin bool) {
	fmt.Println("pfstpaths:", pfstpaths)
	fmt.Println("sigpaths:", sigpaths)
	idxs := psIndexRepsTop(breed, "unbitted")
	frstring := "partialreps"
	if fullreps {
		frstring = "fullreps"
	}
	nowinstring := ""
	if nowin {
		nowinstring = "_nowin"
	}
	fmt.Fprintf(w, `plot_goods_pfst_reps_together_reptop%v \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	%v \
	%v_plfmt.bed \
	good4plots_together_reps_reptop_%v_%v%v.png \
> together_reps_reptop_out_%v_%v%v.txt
`,
		nowinstring,
		pfstpaths[idxs[0]],
		sigpaths[idxs[0]],
		pfstpaths[idxs[1]],
		sigpaths[idxs[1]],
		pfstpaths[idxs[2]],
		sigpaths[idxs[2]],
		pfstpaths[idxs[3]],
		sigpaths[idxs[3]],
		pfstpaths[idxs[4]],
		sigpaths[idxs[4]],
		breed,
		frstring,
		nowinstring,
		breed,
		frstring,
		nowinstring,
	)
}

func FprintPlotTogetherRepsTop(w io.Writer, pfstpaths, sigpaths []string, fullreps, nowin bool) {
	fmt.Fprintf(w, `#!/bin/bash
`)
	for _, breed := range []string{"black", "white", "figurita", "runt"} {
		fmt.Fprintf(w, "\n")
		FprintPlotTogetherRepsTopBreed(w, pfstpaths, sigpaths, breed, fullreps, nowin)
	}
}

func BreedOrder(breed string) int {
	switch breed {
	case "black": return 0
	case "white": return 1
	case "figurita": return 2
	case "runt": return 3
	case "feral": return 4
	default: return 5
	}
	return 5
}

func RepOrder(breed string) int {
	switch breed {
	case "All": return 0
	case "1": return 1
	case "2": return 2
	case "3": return 3
	case "4": return 4
	default: return 5
	}
	return 5
}

func FilterAllrepExp(cfgs []ComboConfig) []ComboConfig {
	var out []ComboConfig
	for _, cfg := range cfgs {
		if cfg.Treatment.Replicate == "All" && cfg.ComparisonType == "experimental" {
			out = append(out, cfg)
		}
	}
	less := func(i, j int) bool {
		return BreedOrder(out[i].Treatment.Breed) < BreedOrder(out[j].Treatment.Breed)
	}
	sort.Slice(out, less)
	return out
}

func FilterEveryRepExp(cfgs []ComboConfig) []ComboConfig {
	var out []ComboConfig
	for _, cfg := range cfgs {
		if cfg.ComparisonType == "experimental" {
			out = append(out, cfg)
		}
	}
	less := func(i, j int) bool {
		if RepOrder(out[i].Treatment.Replicate) < RepOrder(out[j].Treatment.Replicate) {
			return true
		}
		if RepOrder(out[i].Treatment.Replicate) > RepOrder(out[j].Treatment.Replicate) {
			return false
		}
		return BreedOrder(out[i].Treatment.Breed) < BreedOrder(out[j].Treatment.Breed)
	}
	sort.Slice(out, less)
	return out
}

func ReadAllrepExpPlotSets(r io.Reader) []PlotSet {
	cfgs, err := ReadComboConfig(r)
	if err != nil {
		panic(err)
	}
	fcfgs := FilterAllrepExp(cfgs)
	return ConfigsToPlotSets(fcfgs...)
}

func ReadEveryRepExpPlotSets(r io.Reader) []PlotSet {
	cfgs, err := ReadComboConfig(r)
	if err != nil {
		panic(err)
	}
	fcfgs := FilterEveryRepExp(cfgs)
	return ConfigsToPlotSets(fcfgs...)
}

func PlotSetToTogetherSubFullScriptRepsTop() {
	var fullreps bool
	var nowin bool
	flag.BoolVar(&fullreps, "f", false, "use all indivs from each rep")
	flag.BoolVar(&nowin, "n", false, "use unwindowed values")
	flag.Parse()
	plotsets := ReadEveryRepExpPlotSets(os.Stdin)
	fmt.Println("plotsets from ReadEveryRepExpPlotSets", plotsets)
	pfstpaths := []string{}
	sigpaths := []string{}
	for _, set := range plotsets {
		pfstpaths = append(pfstpaths, set.Out + "_pfst_plfmt.bed")
		sigpaths = append(sigpaths, set.GoodPfstSpans)
	}
	FprintPlotTogetherRepsTop(os.Stdout, pfstpaths, sigpaths, fullreps, nowin)
}

// order:
// r1 black unbit
// r1 black bit
// r1 white unbit
// r1 white bit
// r1 fig unbit
// r1 fig bit
// r1 runt unbit
// r1 runt bit
// r2 black unbit
// ...

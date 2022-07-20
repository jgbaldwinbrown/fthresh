package main

import (
	"strings"
	"fmt"
	"io"
	"bufio"
	"os"
	"os/exec"
	"github.com/jgbaldwinbrown/permuvals/permuvals"
)

type PlotSet struct {
	Pfst string
	Fst string
	Selec string
	Out string
}

type PlotSets []PlotSet

func ParsePlotSet(s string) PlotSet {
	line := strings.Split(s, "\t")
	return PlotSet{
		Pfst: line[0],
		Fst: line[1],
		Selec: line[2],
		Out: line[3],
	}
}

func ReadPlotSets(r io.Reader) (ps PlotSets) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		ps = append(ps, ParsePlotSet(s.Text()))
	}
	return
}

func MergeHitsString() string {
	return `#!/bin/bash
set -e

awk -F "\t" -v OFS="\t" '$'${1}' >= '${2}'{$3=sprintf("%d", $3); print $0}' \
> ${3}_thresholded.bed

bedtools merge -i ${3}_thresholded.bed > ${3}_thresh_merge.bed`
}


func ThreshAndMerge(inpath string, col int, thresh float64, outpath string) (err error) {
	in, err := os.Open(inpath)
	if err != nil { return }
	defer in.Close()

	script, err := os.Create("merge_hits_scripted.sh")
	if err != nil { return }
	fmt.Fprintln(script, MergeHitsString())
	script.Close()

	cmd := exec.Command("bash", "merge_hits_scripted.sh", fmt.Sprint(col), fmt.Sprint(thresh), outpath)
	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	return
}

func PfstPrefix(prefix string) string {
	return prefix + "_pfst_plfmt"
}

func FstPrefix(prefix string) string {
	return prefix + "_fst_plfmt"
}

func SelecPrefix(prefix string) string {
	return prefix + "_selec_plfmt_bedified"
}

func MergeOutPrefix(prefix string) string {
	return prefix + "_tm"
}

func FinalOutPath(prefix string) string {
	return prefix + "_tm_probs.txt"
}

func FstMergeOut(prefix string) string {
	return MergeOutPrefix(FstPrefix(prefix)) + "_thresh_merge.bed"
}

func PfstMergeOut(prefix string) string {
	return MergeOutPrefix(PfstPrefix(prefix)) + "_thresh_merge.bed"
}

func SelecMergeOut(prefix string) string {
	return MergeOutPrefix(SelecPrefix(prefix)) + "_thresh_merge.bed"
}

func main() {
	plot_sets := ReadPlotSets(os.Stdin)
	seed := 0
	for _, set := range plot_sets {
		err := ThreshAndMerge(PfstPrefix(set.Out)+".bed", 9, 1000.0, MergeOutPrefix(PfstPrefix(set.Out)))
		if err != nil { panic(err) }
		err = ThreshAndMerge(FstPrefix(set.Out) + ".bed", 4, .10, MergeOutPrefix(FstPrefix(set.Out)))
		if err != nil { panic(err) }
		err = ThreshAndMerge(SelecPrefix(set.Out) + ".bed", 4, .04, MergeOutPrefix(SelecPrefix(set.Out)))
		if err != nil { panic(err) }

		bedpaths := []string{PfstMergeOut(set.Out), FstMergeOut(set.Out), SelecMergeOut(set.Out)}
		// bedpaths := []string{MergeOutPrefix(set.Out + "_pfst"), MergeOutPrefix(set.Out + "_fst"), MergeOutPrefix(set.Out + "_selec")}
		bedpaths_conn, err := os.Create("bedpaths_scripted.txt")
		if err != nil { panic(err) }
		for _, b := range bedpaths {
			fmt.Fprintln(bedpaths_conn, b)
		}
		bedpaths_conn.Close()

		var flags permuvals.Flags
		flags.BedPaths = "bedpaths_scripted.txt"
		flags.GenomeBedPath = "louse_genome_0.1.1.chrlens.bed"
		flags.Iterations = 2000
		flags.Rseed = seed
		fmt.Println("start comp")
		comp, err := permuvals.FullCompare(flags)
		fmt.Println("end comp")
		if err != nil { panic(err) }

		outconn, err := os.Create(FinalOutPath(set.Out))
		if err != nil { panic(err) }
		permuvals.FprintProbs(outconn, comp.Probs)
		outconn.Close()

		seed++
	}

}

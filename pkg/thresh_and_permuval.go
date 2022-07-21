package fthresh

import (
	"strings"
	"fmt"
	"io"
	"bufio"
	"os"
	"os/exec"
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

////////////////////////////////////////////////////////////////////////////////////////////////////
func PfstOut(prefix string) string {
	return PfstPrefix(prefix) + ".bed"
}

func GggenesOut(prefix, chr, start, end string) string {
	return prefix + "_" + chr + "_" + start + "_" + end + ".pdf"
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func RunThreshMergeAndPermuval() {
	plot_sets := ReadPlotSets(os.Stdin)
	seed := 0
	for _, set := range plot_sets {
		err := ThreshMergeAndPerm(set, "nostat", seed)
		if err != nil {
			panic(err)
		}
		seed++
	}

}

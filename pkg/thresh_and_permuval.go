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
	GoodPfstSpans string
}

type PlotSets []PlotSet

func ParsePlotSet(s string) PlotSet {
	line := strings.Split(s, "\t")
	set := PlotSet{
		Pfst: line[0],
		Fst: line[1],
		Selec: line[2],
		Out: line[3],
	}
	if len(line) >= 5 {
		set.GoodPfstSpans = line[4]
	}
	return set
}

func ConfigToPlotSet(c ComboConfig) PlotSet {
	out := PlotSet{}
	out.Pfst = c.Pfst
	out.Fst = c.Fst
	out.Selec = c.Selec
	out.Out = c.OutPrefix
	out.GoodPfstSpans = c.Subtractions
	return out
}

func ConfigsToPlotSets(cs ...ComboConfig) []PlotSet {
	out := make([]PlotSet, len(cs))
	for i, c := range cs {
		out[i] = ConfigToPlotSet(c)
	}
	return out
}

func ReadCfgsAndPlotSets(r io.Reader) ([]ComboConfig, PlotSets) {
	cfgs, err := ReadComboConfig(r)
	if err != nil {
		panic(err)
	}
	return cfgs, ConfigsToPlotSets(cfgs...)
}
func ReadCfgPlotSets(r io.Reader) PlotSets {
	_, plotsets := ReadCfgsAndPlotSets(r)
	return plotsets
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

mawk -F "\t" -v OFS="\t" '$'${1}' >= '${2}'{
	$3=sprintf("%d", $3);
	if ($2 < 0) { $2 = 0 };
	if ($3 < 0) { $3 = 0 };
	print $0
}' \
> ${3}_thresholded.bed

bedtools merge -i ${3}_thresholded.bed > ${3}_thresh_merge.bed`
}

func WriteFile(path, content string) error {
	script, err := os.Create(path)
	if err != nil {
		return err
	}
	fmt.Fprintln(script, content)
	script.Close()
	return nil
}

func ThreshAndMerge(inpath string, col int, thresh float64, outpath string) (err error) {
	in, err := os.Open(inpath)
	if err != nil { return }
	defer in.Close()

	err = WriteFile("merge_hits_scripted.sh", MergeHitsString())
	if err != nil {
		return
	}

	fmt.Println("col:", col)
	fmt.Println("thresh:", thresh)
	fmt.Println("outpath:", outpath)
	cmd := exec.Command("bash", "merge_hits_scripted.sh", fmt.Sprint(col), fmt.Sprint(thresh), outpath)
	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	return err
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

func FstPercMergeOut(prefix string) string {
	return MergeOutPrefix(FstPrefix(prefix)) + "_perc_thresh_merge.bed"
}

func PfstPercMergeOut(prefix string) string {
	return MergeOutPrefix(PfstPrefix(prefix)) + "_perc_thresh_merge.bed"
}

func SelecPercMergeOut(prefix string) string {
	return MergeOutPrefix(SelecPrefix(prefix)) + "_perc_thresh_merge.bed"
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func PfstOut(prefix string) string {
	return PfstPrefix(prefix) + ".bed"
}

func GggenesOut(prefix, chr, start, end string) string {
	return prefix + "_" + chr + "_" + start + "_" + end + ".pdf"
}

func GggenesOutNofstNoselec(prefix, chr, start, end string) string {
	return prefix + "_" + chr + "_" + start + "_" + end + "_nofst_noselec.pdf"
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func RunThreshMergeAndPermuval() {
	plot_sets := ReadPlotSets(os.Stdin)
	seed := 0
	for _, set := range plot_sets {
		err := ThreshMergeAndPerm(set, "pfst_allcombo", true, seed)
		if err != nil {
			panic(err)
		}
		seed++
	}

}

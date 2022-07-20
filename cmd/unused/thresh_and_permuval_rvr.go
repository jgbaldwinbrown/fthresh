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

func ThreshMergeAll(set PlotSet) (err error) {
	err = ThreshAndMerge(PfstPrefix(set.Out)+".bed", 9, 1000.0, MergeOutPrefix(PfstPrefix(set.Out)))
	if err != nil { return }
	err = ThreshAndMerge(FstPrefix(set.Out) + ".bed", 4, .10, MergeOutPrefix(FstPrefix(set.Out)))
	if err != nil { return }
	err = ThreshAndMerge(SelecPrefix(set.Out) + ".bed", 4, .04, MergeOutPrefix(SelecPrefix(set.Out)))
	if err != nil { return }

	return nil
}

func TMAParallel(set PlotSet, errchan chan error) {
	err := ThreshMergeAll(set)
	errchan <- err
}

func PermAll(bedpaths []string, statistic string, seed int) (err error) {
	bedpaths_path := "bedpaths_scripted_" + statistic + ".txt"
	bedpaths_conn, err := os.Create(bedpaths_path)
	if err != nil { return }
	for _, b := range bedpaths {
		fmt.Fprintln(bedpaths_conn, b)
	}
	bedpaths_conn.Close()

	var flags permuvals.Flags
	flags.BedPaths = bedpaths_path
	flags.GenomeBedPath = "louse_genome_0.1.1.chrlens.bed"
	flags.Iterations = 2000
	flags.Rseed = seed
	fmt.Println("start comp")
	comp, err := permuvals.FullCompare(flags)
	fmt.Println("end comp")
	if err != nil { return }

	outconn, err := os.Create(statistic + "_probs.txt")
	if err != nil { return }
	permuvals.FprintProbs(outconn, comp.Probs)
	outconn.Close()

	return nil
}

func PAParallel(bedpaths []string, statistic string, seed int, errchan chan error) {
	err := PermAll(bedpaths, statistic, seed)
	errchan <- err
}

func TMASets(sets PlotSets) (errors []error) {
	var errorchans []chan error
	for _, set := range sets {
		ec := make(chan error)
		errorchans = append(errorchans, ec)
		go TMAParallel(set, ec)
	}
	for _, ec := range errorchans {
		errors = append(errors, <-ec)
	}
	return
}

func PASets(all_bedpaths [][]string, statistics []string) (errors []error) {
	var errorchans []chan error
	for seed, bedpaths := range all_bedpaths {
		ec := make(chan error)
		errorchans = append(errorchans, ec)
		go PAParallel(bedpaths, statistics[seed], seed, ec)
	}
	for _, ec := range errorchans {
		errors = append(errors, <-ec)
	}
	return
}

func GetAllBedpaths(sets PlotSets) (out [][]string) {
	out = append(out, []string{}, []string{}, []string{})
	for _, set := range sets {
		out[0] = append(out[0], PfstMergeOut(set.Out))
		out[1] = append(out[1], FstMergeOut(set.Out))
		out[2] = append(out[2], SelecMergeOut(set.Out))
	}
	return out
}

func main() {
	plot_sets := ReadPlotSets(os.Stdin)

	errors := TMASets(plot_sets)
	for _, err := range errors {
		if err != nil { panic(err) }
	}

	all_bedpaths := GetAllBedpaths(plot_sets)
	statistics := []string{"pFst", "Fst", "Selec"}
	errors = PASets(all_bedpaths, statistics)
	for _, err := range errors {
		if err != nil { panic(err) }
	}
}

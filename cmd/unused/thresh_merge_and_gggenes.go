package main

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

type GggenesArgs struct {
	annotpath string
	datapath string
	snppath string
	outpath string
	chr string
	start string
	end string
	height string
	width string
}

func DefaultGargs() GggenesArgs {
	a := GggenesArgs{}
	a.outpath = "out.pdf"
	a.height = "3"
	a.width = "8"
	return a
}

func Gggenes(a GggenesArgs) error {
	return GggenesInternal(a.annotpath, a.datapath, a.snppath, a.outpath, a.chr, a.start, a.end, a.height, a.width)
}

func GggenesInternal(annotpath, datapath, snppath, outpath, chr, start, end, height, width string) error {
	cmd := exec.Command(
		"Rscript",
		"plot_region.R",
		annotpath,
		datapath,
		snppath,
		fmt.Sprintf("%v:%v:%v", chr, start, end),
		outpath,
		fmt.Sprintf("%v:%v", width, height),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GggenesWorker(inputs <-chan GggenesArgs, results chan<- error) {
	for a := range inputs {
		results <- Gggenes(a)
	}
}

func GggenesBed(bed [][]string, a GggenesArgs, outprefix string) error {
	numJobs := len(bed)
	jobs := make(chan GggenesArgs, numJobs)
	errs := make(chan error, numJobs)

	for i:=0; i<8; i++ {
		go GggenesWorker(jobs, errs)
	}

	for _, entry := range bed {
		a2 := a
		a2.chr = entry[0]
		a2.start = entry[1]
		a2.end = entry[2]
		a2.outpath = GggenesOut(outprefix, a2.chr, a2.start, a2.end)

		jobs <- a2
	}
	close(jobs)

	for err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func ReadTable(path string) ([][]string, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	out := [][]string{}
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	for s.Scan() {
		out = append(out, strings.Split(s.Text(), "\t"))
	}
	return out, nil
}

// # $1: input selection coefficients
// # $2: input fsts
// # $3: input pfsts
// # $4: output combined data bed
func MakeDataCombo(pfstpath, fstpath, selecpath, outpre string) (datapath string, err error) {
	datapath = outpre + "_databed.bed"
	fmt.Println("data: ", selecpath, fstpath, pfstpath, datapath)
	cmd := exec.Command("./setup3.sh", selecpath, fstpath, pfstpath, datapath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return datapath, err
}

func MakeSnpBed(snpvcfpath string) (snpbedpath string, err error) {
	snpbedpath = snpvcfpath + "_snpbed.bed"
	fmt.Println(snpvcfpath, snpbedpath)
	cmd := exec.Command("./setup3snp.sh", snpvcfpath, snpbedpath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return snpbedpath, err
}

func WriteBed(bed [][]string, path string) error {
	outconn, err := os.Open(path)
	if err != nil {
		return err
	}
	defer outconn.Close()

	w := bufio.NewWriter(outconn)
	defer w.Flush()

	for _, b := range bed {
		io.WriteString(w, strings.Join(b, "\t") + "\n")
	}

	return nil
}

func GggenesAllBedpaths(paths_sets [][]string, annotpath, snpvcfpath string) error {
	// stats := []string{"pfst", "fst", "selec", "outpre"}
	snpbedpath, err := MakeSnpBed(snpvcfpath)
	if err != nil {
		return err
	}

	for i, _ := range paths_sets[0] {
		datapath, err := MakeDataCombo(paths_sets[4][i], paths_sets[1][i], paths_sets[2][i], paths_sets[3][i])
		if err != nil {
			return err
		}

		bed, err := ReadTable(paths_sets[0][i])
		if err != nil {
			return err
		}

		args := DefaultGargs()
		args.annotpath = annotpath
		args.datapath = datapath
		args.snppath = snpbedpath

		outpre := paths_sets[3][i]
		_ = WriteBed(bed, outpre + "_thresholded_intervals.bed")
		GggenesBed(bed, args, outpre + "_gggenes")
	}
	return nil
}

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

awk -F "\t" -v OFS="\t" '$'${1}' >= '${2}'{$2=sprintf("%d", $2-25000); if ($2 < 0){$2=0}; $3=sprintf("%d", $3+25000); print $0}' \
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
	return prefix + "_threshmerge2"
}

func FstMergeOut(prefix string) string {
	return FstPrefix(prefix) + ".bed"
}

func PfstMergeOut(prefix string) string {
	return MergeOutPrefix(PfstPrefix(prefix)) + "_thresh_merge.bed"
}

func PfstOut(prefix string) string {
	return PfstPrefix(prefix) + ".bed"
}

func SelecMergeOut(prefix string) string {
	return SelecPrefix(prefix) + ".bed"
}

func GggenesOut(prefix, chr, start, end string) string {
	return prefix + "_" + chr + "_" + start + "_" + end + ".pdf"
}

func ThreshMergeAll(set PlotSet) (err error) {
	err = ThreshAndMerge(PfstPrefix(set.Out)+".bed", 9, 25, MergeOutPrefix(PfstPrefix(set.Out)))
	if err != nil { return }
	err = ThreshAndMerge(FstPrefix(set.Out) + ".bed", 4, .19, MergeOutPrefix(FstPrefix(set.Out)))
	if err != nil { return }
	err = ThreshAndMerge(SelecPrefix(set.Out) + ".bed", 4, .03, MergeOutPrefix(SelecPrefix(set.Out)))
	if err != nil { return }

	return nil
}

func TMAParallel(set PlotSet, errchan chan error) {
	err := ThreshMergeAll(set)
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

func GetAllBedpaths(sets PlotSets) (out [][]string) {
	out = append(out, []string{}, []string{}, []string{}, []string{}, []string{})
	for _, set := range sets {
		out[0] = append(out[0], PfstMergeOut(set.Out))
		out[1] = append(out[1], FstMergeOut(set.Out))
		out[2] = append(out[2], SelecMergeOut(set.Out))
		out[3] = append(out[3], set.Out + "_gggenes")
		out[4] = append(out[4], PfstOut(set.Out))
	}
	return out
}

func main() {
	plot_sets := ReadPlotSets(os.Stdin)

	errors := TMASets(plot_sets)
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}

	all_bedpaths := GetAllBedpaths(plot_sets)

	annotpath := "louse_annotation_0.1.1.gff.gz"
	snppath := "snpdat_out_forplot.txt.gz"
	err := GggenesAllBedpaths(all_bedpaths, annotpath, snppath)
	if err != nil {
		panic(err)
	}
}

package fthresh

import (
	"strings"
	"fmt"
	"io"
	"bufio"
	"os"
	"os/exec"
	"github.com/jgbaldwinbrown/pmap/pkg"
)

const setupShTxt string = `#!/bin/bash
set -e

# $1: input selection coefficients
# $2: input fsts
# $3: input pfsts
# $4: output combined data bed

cat "$1" | \
awkf '{print $1, $2, $3, $4, "selec"}' \
> "$4"

cat "$2" | \
awkf '{print $1, $2, $3, $4, "fst"}' \
>> "$4"

cat "$3" | \
awkf '{print $1, $2, $3, $9, "pfst"}' \
>> "$4"
`

const setupSnpShTxt string = `#!/bin/bash
set -e

# $1: input snps vcf
# $2: output edited snps bed

gunzip -c "$1" | \
awkf '$1!~/^#/ {print $1, $2-1, $2, "snp"}' \
> "$2"
`

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

func GggenesBed(bed [][]string, a GggenesArgs, outprefix string) error {
	var inputs []GggenesArgs
	for _, entry := range bed {
		a2 := a
		a2.chr = entry[0]
		a2.start = entry[1]
		a2.end = entry[2]
		a2.outpath = GggenesOut(outprefix, a2.chr, a2.start, a2.end)
		inputs = append(inputs, a2)
	}
	errs := pmap.Map(Gggenes, inputs, 2)
	for _, err := range errs {
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
	err = WriteFile("setup.sh", setupShTxt)
	if err != nil { return }

	datapath = outpre + "_databed.bed"
	fmt.Println("data: ", selecpath, fstpath, pfstpath, datapath)
	cmd := exec.Command("sh", "setup.sh", selecpath, fstpath, pfstpath, datapath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return datapath, err
}

func MakeSnpBed(snpvcfpath string) (snpbedpath string, err error) {
	err = WriteFile("setup_snp.sh", setupSnpShTxt)
	if err != nil { return }

	snpbedpath = snpvcfpath + "_snpbed.bed"
	fmt.Println(snpvcfpath, snpbedpath)
	cmd := exec.Command("sh", "setup_snp.sh", snpvcfpath, snpbedpath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return snpbedpath, err
}

func WriteBed(bed [][]string, path string) error {
	outconn, err := os.Create(path)
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


func GetAllBedpathsGggenes(sets PlotSets) (out [][]string) {
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

func ThreshMergeAndGggenesMini() {
	annotpath := "louse_annotation_0.1.1_chr1.gff.gz"
	snpbedpath := "snpdat_out_forplot.txt.gz_snpbed.bed"

	args := DefaultGargs()
	args.annotpath = annotpath
	args.datapath = "_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Low_High___multiplot_gggenes_databed.bed"
	args.snppath = snpbedpath

	bed, err := ReadTable("chr1mini_thresh_merge.bed")
	if err != nil {
		panic(err)
	}

	prefix := "minigggenes_out2"
	err = WriteBed(bed, prefix + "_thresholded_intervals.bed")
	if err != nil {
		panic(err)
	}
	err = GggenesBed(bed, args, prefix)
	if err != nil {
		panic(err)
	}
}

func ThreshMergeAndGggenes() {
	plot_sets := ReadPlotSets(os.Stdin)

	errors := TMASets(plot_sets)
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}

	all_bedpaths := GetAllBedpathsGggenes(plot_sets)

	annotpath := "louse_annotation_0.1.1.gff.gz"
	snppath := "snpdat_out_forplot.txt.gz"
	err := GggenesAllBedpaths(all_bedpaths, annotpath, snppath)
	if err != nil {
		panic(err)
	}
}

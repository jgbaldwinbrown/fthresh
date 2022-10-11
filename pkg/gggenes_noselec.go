package fthresh

import (
	"fmt"
	"os"
	"os/exec"
	"github.com/jgbaldwinbrown/pmap/pkg"
)

const setupShTxtNofstNoselec string = `#!/bin/bash
set -e

# $1: input selection coefficients
# $2: input fsts
# $3: input pfsts
# $4: output combined data bed

cat "$3" | \
mawk -F "\t" -v OFS="\t" '{print $1, $2, $3, $9, "pfst"}' \
>> "$4"
`

func GggenesNofstNoselec(a GggenesArgs) error {
	cmd := exec.Command(
		"plot_region",
		a.annotpath,
		a.datapath,
		a.snppath,
		fmt.Sprintf("%v:%v:%v", a.chr, a.start, a.end),
		a.outpath,
		fmt.Sprintf("%v:%v", a.width, a.height),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// # $1: input selection coefficients
// # $2: input fsts
// # $3: input pfsts
// # $4: output combined data bed
func MakeDataComboNofstNoselec(pfstpath, fstpath, selecpath, outpre string) (datapath string, err error) {
	return makeDataComboInternal(pfstpath, fstpath, selecpath, outpre, setupShTxtNofstNoselec)
}

func GggenesAllBedpathsNofstNoselec(paths_sets [][]string, annotpath, snpvcfpath string) error {
	// stats := []string{"pfst", "fst", "selec", "outpre"}
	snpbedpath, err := MakeSnpBed(snpvcfpath)
	if err != nil {
		return err
	}

	var errs errList
	for i, _ := range paths_sets[0] {
		datapath, err := MakeDataComboNofstNoselec(paths_sets[4][i], paths_sets[1][i], paths_sets[2][i], paths_sets[3][i])
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
		err = GggenesBedNofstNoselec(bed, args, outpre + "_gggenes")
		if err != nil {
			errs = append(errs, err)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}

func GggenesBedNofstNoselec(bed [][]string, a GggenesArgs, outprefix string) error {
	var inputs []GggenesArgs
	for _, entry := range bed {
		a2 := a
		a2.chr = entry[0]
		a2.start = entry[1]
		a2.end = entry[2]
		a2.outpath = GggenesOutNofstNoselec(outprefix, a2.chr, a2.start, a2.end)
		inputs = append(inputs, a2)
	}
	errs := pmap.Map(GggenesNofstNoselec, inputs, 2)
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}


func GggenesFullNofstNoselec() {
	plot_sets := ReadCfgPlotSets(os.Stdin)

	all_bedpaths := GetAllBedpathsGggenesPerc(plot_sets)

	annotpath := "louse_annotation_0.1.1.gff.gz"
	snppath := "snpdat_out_forplot.txt.gz"
	err := GggenesAllBedpathsNofstNoselec(all_bedpaths, annotpath, snppath)
	if err != nil {
		panic(err)
	}
}

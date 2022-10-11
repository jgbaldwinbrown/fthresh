package fthresh

import (
	"os/exec"
	"fmt"
	// "sync"
	"runtime"
	"strconv"
	"flag"
	"os"
	"io"
	"github.com/jgbaldwinbrown/permuvals/pkg"
	"github.com/jgbaldwinbrown/pmap/pkg"
	"github.com/jgbaldwinbrown/accel/accel"
)

func ExpandMergeString() string {
	return `#!/bin/bash
set -e

mawk -F "\t" -v OFS="\t" '{
	$2=$2 - '${1}';
	$3=$3 + '${1}';
	if ($2 < 0) { $2 = 0 };
	if ($3 < 0) { $3 = 0 };
	print $0;
}' \
> ${2}_thresholded.bed

bedtools merge -i ${2}_thresholded.bed > ${2}_thresh_merge.bed`
}

// func Quantile(data []float64, quantile float64) (threshpos int, thresh float64, overthresh []float64) {

func ExpandAndMerge(inpath string, slop int, outpath string) (err error) {
	in, err := os.Open(inpath)
	if err != nil { return }
	defer in.Close()

	err = WriteFile("expand_merge_scripted.sh", ExpandMergeString())
	if err != nil {
		return
	}

	fmt.Println("slop:", slop)
	fmt.Println("outpath:", outpath)
	cmd := exec.Command("bash", "expand_merge_scripted.sh", fmt.Sprint(slop), outpath)
	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	return err
}

func PercThreshAndMerge(inpath string, col int, percentile float64, slop int, outpath string) error {
	fmt.Println("paths:", inpath, outpath)
	datatxt, err := ReadTable(inpath)
	if err != nil { return err }
	var data []float64
	for _, line := range datatxt {
		var float float64
		float, err = strconv.ParseFloat(line[col-1], 64)
		if err == nil {
			data = append(data, float)
		}
	}
	_, thresh, _ := accel.Quantile(data, percentile)
	// threshpos, thresh, overthresh := accel.Quantile(data, percentile)
	// fmt.Println("threshold stuff:")
	// fmt.Println(data[:10], percentile, threshpos, thresh, overthresh)
	// fmt.Println("paths again:", inpath, outpath)
	err = ThreshAndMerge(inpath, col, thresh, outpath + "_unexpanded")
	if err != nil {
		return err
	}

	err = ExpandAndMerge(outpath + "_unexpanded_thresh_merge.bed", slop, outpath);
	if err != nil {
		return err
	}
	fmt.Println("paths again2:", inpath, outpath)
	return nil
}

func PercThreshMergeAll(set PlotSet, percentile float64) (err error) {
	err = PercThreshAndMerge(PfstPrefix(set.Out)+".bed", 9, percentile, 250000, MergeOutPrefix(PfstPrefix(set.Out)) + "_perc")
	if err != nil { return }
	err = PercThreshAndMerge(FstPrefix(set.Out) + ".bed", 4, percentile, 250000, MergeOutPrefix(FstPrefix(set.Out)) + "_perc")
	if err != nil { return }
	// err = PercThreshAndMerge(SelecPrefix(set.Out) + ".bed", 4, percentile, 250000, MergeOutPrefix(SelecPrefix(set.Out)) + "_perc")
	// if err != nil { return }

	return nil
}

type Comp struct {
	Breed1 string
	Bit1 string
	Rep1 int
	Gen1 int
	Breed2 string
	Bit2 string
	Rep2 int
	Gen2 int
}

func (c Comp) Path(statistic string) string {
	return BedString(c.Breed1, c.Bit1, c.Rep1, c.Gen1, c.Breed2, c.Bit2, c.Rep2, c.Gen2, statistic)
}

func (c Comp) PathFull(statistic string) string {
	return BedStringFull(c.Breed1, c.Bit1, c.Rep1, c.Gen1, c.Breed2, c.Bit2, c.Rep2, c.Gen2, statistic)
}

func (c Comp) OutputPrefix(statistic string) string {
	return CompOutput(c.Breed1, c.Bit1, statistic, c.Rep1)
}

type GoodAndAlts struct {
	Good Comp
	Alts []Comp
}

func PrebuiltGoodAndAlts() []GoodAndAlts {
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
		GoodAndAlts{
			Good: Comp{"figurita", "unbitted", 0, 36, "feral", "unbitted", 0, 36},
			Alts: []Comp {
				Comp{"figurita", "bitted", 0, 36, "feral", "bitted", 0, 36},
			},
		},
		GoodAndAlts{
			Good: Comp{"runt", "unbitted", 0, 36, "feral", "unbitted", 0, 36},
			Alts: []Comp {
				Comp{"runt", "bitted", 0, 36, "feral", "bitted", 0, 36},
			},
		},
	}
}

func ReadPath[T any](path string, f func(io.Reader) (T, error)) (T, error) {
	var t T
	r, err := os.Open(path)
	if err != nil {
		return t, err
	}
	defer r.Close()
	t, err = f(r)
	return t, err
}

type Plot4Flags struct {
	Threads int
	Percentile float64
	GoodsAndAltsPath string
	SubFull bool
	FullReps bool
}

func GetPlot4Flags() Plot4Flags {
	var f Plot4Flags

	flag.IntVar(&f.Threads, "t", 1, "Threads to use")
	percstring := flag.String("p", ".001", "Percentile threshold")
	flag.BoolVar(&f.SubFull, "s", false, "Set to subtract entire region if it partially intersects with alt")
	flag.BoolVar(&f.FullReps, "f", false, "Use all indivs in each replicate, not just top/bottom 20%")
	// flag.StringVar(&f.GoodsAndAltsPath, "g", "", "Path to tab-separated sets of good and alt comparisons")
	flag.Parse()

	var err error
	f.Percentile, err = strconv.ParseFloat(*percstring, 64)
	if err != nil {
		panic(err)
	}
	// if f.GoodsAndAltsPath == "" {
	// 	panic(fmt.Errorf("missing GoodsAndAltsPath (-g)"))
	// }
	return f
}

func SubtractAlts(gset GoodAndAlts, statistic string, subfull bool, fullreps bool) (outpath string, err error) {
	fmt.Fprintln(os.Stderr, "subtracting:", gset, statistic, subfull)
	goodpath := gset.Good.Path(statistic)
	if fullreps {
		goodpath = gset.Good.PathFull(statistic)
	}
	good, err := ReadPath(goodpath, func(r io.Reader) (permuvals.Bed, error) {
		return permuvals.GetBed(r, goodpath)
	})
	if err != nil {
		fmt.Fprintln(os.Stderr,"Error!")
		return "", err
	}
	for _ , altcomp := range gset.Alts {
		path := altcomp.Path(statistic)
		if fullreps {
			path = altcomp.PathFull(statistic)
		}
		fmt.Println("path to subtract:", path)
		alt, err := ReadPath(path, func(r io.Reader) (permuvals.Bed, error) {
			return permuvals.GetBed(r, path)
		})
		if err != nil {
			fmt.Fprintln(os.Stderr,"Error2!")
			return "", err
		}
		if subfull {
			good = good.SubtractFullsBed(alt)
		} else {
			good.SubtractBed(alt)
		}
	}

	if subfull {
		outpath = gset.Good.OutputPrefix(statistic) + "_subfulls.bed"
	} else {
		outpath = gset.Good.OutputPrefix(statistic) + ".bed"
	}
	ovlconn, err := os.Create(outpath)
	if err != nil {
		fmt.Fprintln(os.Stderr,"Error3!")
		return "", err
	}
	defer ovlconn.Close()

	permuvals.FprintBeds(ovlconn, good)

	return outpath, err
}

type subtractAllArgs struct {
	gset GoodAndAlts
	statistic string
	subfull bool
	fullreps bool
}

type subtractAllOuts struct {
	outpath string
	err error
}

func aggregateSubtractArgs(gsets []GoodAndAlts, statistics []string, subfull, fullreps bool) []subtractAllArgs {
	njobs := len(gsets) * len(statistics)
	mod := len(statistics)
	out := make([]subtractAllArgs, njobs)
	for job, _ := range out {
		out[job].gset = gsets[job/mod]
		out[job].statistic = statistics[job%mod]
		out[job].subfull = subfull
	}
	return out
}

func SubtractAllAlts(gsets []GoodAndAlts, statistics []string, subfull, fullreps bool, threads int) (outpaths []string, errs []error) {
	args := aggregateSubtractArgs(gsets, statistics, subfull, fullreps)
	f := func(a subtractAllArgs) subtractAllOuts {
		var o subtractAllOuts
		o.outpath, o.err = SubtractAlts(a.gset, a.statistic, a.subfull, a.fullreps)
		return o
	}
	out := pmap.Map(f, args, threads)
	outpaths = make([]string, len(out))
	errs = make([]error, len(out))
	for i, val := range out {
		outpaths[i] = val.outpath
		errs[i] = val.err
	}
	return outpaths, errs
}

func PercTMASets(sets PlotSets, percentile float64) (errors []error) {
	f := func(set PlotSet) error {
		return PercThreshMergeAll(set, percentile)
	}
	return pmap.Map(f, sets, -1)
}

func RunGood4Plots() {
	flags := GetPlot4Flags()
	// plot_sets := ReadPlotSets(os.Stdin)
	cfgs, err := ReadComboConfig(os.Stdin)
	if err != nil {
		panic(err)
	}
	plot_sets := ConfigsToPlotSets(cfgs...)
	runtime.GOMAXPROCS(flags.Threads)

	errors := PercTMASets(plot_sets, flags.Percentile)
	for _, err := range errors {
		if err != nil { panic(err) }
	}

	goodsAndAlts := PrebuiltGoodAndAlts()
	statistics := []string{"pFst", "Fst", "Selec"}
	outpaths, errors := SubtractAllAlts(goodsAndAlts, statistics, false, flags.FullReps, flags.Threads)
	for _, err := range errors {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}

	pathsconn, err := os.Create("goodbedpaths.txt")
	if err != nil { panic(err) }
	defer pathsconn.Close()
	for _, path := range outpaths {
		fmt.Fprintln(pathsconn, path)
	}

	outpaths_f, errors_f := SubtractAllAlts(goodsAndAlts, statistics, true, flags.FullReps, flags.Threads)
	for _, err := range errors_f {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}

	pathsconn_f, err_f := os.Create("goodbedpaths2_f.txt")
	if err_f != nil { panic(err_f) }
	defer pathsconn_f.Close()
	for _, path := range outpaths_f {
		fmt.Fprintln(pathsconn_f, path)
	}
}

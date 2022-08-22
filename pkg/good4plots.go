package fthresh

import (
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


// func Quantile(data []float64, quantile float64) (threshpos int, thresh float64, overthresh []float64) {

func PercThreshAndMerge(inpath string, col int, percentile float64, outpath string) error {
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
	threshpos, thresh, overthresh := accel.Quantile(data, percentile)
	fmt.Println("threshold stuff:")
	fmt.Println(data[:10], percentile, threshpos, thresh, overthresh)
	fmt.Println("paths again:", inpath, outpath)
	err = ThreshAndMerge(inpath, col, thresh, outpath)
	if err != nil {
		return err
	}
	fmt.Println("paths again2:", inpath, outpath)
	return nil
}

func PercThreshMergeAll(set PlotSet, percentile float64) (err error) {
	err = PercThreshAndMerge(PfstPrefix(set.Out)+".bed", 9, percentile, MergeOutPrefix(PfstPrefix(set.Out)) + "_perc")
	if err != nil { return }
	err = PercThreshAndMerge(FstPrefix(set.Out) + ".bed", 4, percentile, MergeOutPrefix(FstPrefix(set.Out)) + "_perc")
	if err != nil { return }
	err = PercThreshAndMerge(SelecPrefix(set.Out) + ".bed", 4, percentile, MergeOutPrefix(SelecPrefix(set.Out)) + "_perc")
	if err != nil { return }

	return nil
}

type Comp struct {
	Breed1 string
	Bit1 string
	Breed2 string
	Bit2 string
}

func (c Comp) Path(statistic string) string {
	return BedString(c.Breed1, c.Bit1, c.Breed2, c.Bit2, statistic)
}

func (c Comp) OutputPrefix(statistic string) string {
	return CompOutput(c.Breed1, c.Bit1, statistic)
}

type GoodAndAlts struct {
	Good Comp
	Alts []Comp
}

func PrebuiltGoodAndAlts() []GoodAndAlts {
	return []GoodAndAlts {
		GoodAndAlts{
			Good: Comp{"black", "unbitted", "feral", "unbitted"},
			Alts: []Comp {
				Comp{"black", "bitted", "feral", "bitted"},
				Comp{"black", "bitted", "white", "bitted"},
			},
		},
		GoodAndAlts{
			Good: Comp{"white", "unbitted", "feral", "unbitted"},
			Alts: []Comp {
				Comp{"white", "bitted", "feral", "bitted"},
				Comp{"white", "bitted", "black", "bitted"},
			},
		},
		GoodAndAlts{
			Good: Comp{"figurita", "unbitted", "feral", "unbitted"},
			Alts: []Comp {
				Comp{"figurita", "bitted", "feral", "bitted"},
				Comp{"figurita", "bitted", "runt", "bitted"},
			},
		},
		GoodAndAlts{
			Good: Comp{"runt", "unbitted", "feral", "unbitted"},
			Alts: []Comp {
				Comp{"runt", "bitted", "feral", "bitted"},
				Comp{"runt", "bitted", "figurita", "bitted"},
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
}

func GetPlot4Flags() Plot4Flags {
	var f Plot4Flags

	flag.IntVar(&f.Threads, "t", 1, "Threads to use")
	percstring := flag.String("p", ".001", "Percentile threshold")
	flag.BoolVar(&f.SubFull, "s", false, "Set to subtract entire region if it partially intersects with alt")
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

func SubtractAlts(gset GoodAndAlts, statistic string, subfull bool) (outpath string, err error) {
	good, err := ReadPath(gset.Good.Path(statistic), func(r io.Reader) (permuvals.Bed, error) {
		return permuvals.GetBed(r, gset.Good.Path(statistic))
	})
	if err != nil {
		return "", err
	}
	for _ , altcomp := range gset.Alts {
		path := altcomp.Path(statistic)
		alt, err := ReadPath(path, func(r io.Reader) (permuvals.Bed, error) {
			return permuvals.GetBed(r, path)
		})
		if err != nil {
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
}

type subtractAllOuts struct {
	outpath string
	err error
}

func aggregateSubtractArgs(gsets []GoodAndAlts, statistics []string, subfull bool) []subtractAllArgs {
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

func SubtractAllAlts(gsets []GoodAndAlts, statistics []string, subfull bool, threads int) (outpaths []string, errs []error) {
	args := aggregateSubtractArgs(gsets, statistics, subfull)
	f := func(a subtractAllArgs) subtractAllOuts {
		var o subtractAllOuts
		o.outpath, o.err = SubtractAlts(a.gset, a.statistic, a.subfull)
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
	plot_sets := ReadPlotSets(os.Stdin)
	runtime.GOMAXPROCS(flags.Threads)

	errors := PercTMASets(plot_sets, flags.Percentile)
	for _, err := range errors {
		if err != nil { panic(err) }
	}

	goodsAndAlts := PrebuiltGoodAndAlts()
	statistics := []string{"pFst", "Fst", "Selec"}
	outpaths, errors := SubtractAllAlts(goodsAndAlts, statistics, flags.SubFull, flags.Threads)
	for _, err := range errors {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}

	pathsconn, err := os.Create("goodbedpaths.txt")
	if err != nil { panic(err) }
	defer pathsconn.Close()
	for _, path := range outpaths {
		fmt.Fprintln(pathsconn, path)
	}
}

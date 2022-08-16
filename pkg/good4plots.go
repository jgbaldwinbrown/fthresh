package fthresh

import (
	"fmt"
	"sync"
	"runtime"
	"strconv"
	"flag"
	"os"
	// "golang.org/x/exp/slices"
	// "github.com/jgbaldwinbrown/lscan/pkg"
	"io"
	"github.com/jgbaldwinbrown/permuvals/pkg"
)

func PercThreshMergeAll(set PlotSet, percentile float64) (err error) {
	err = ThreshAndMerge(PfstPrefix(set.Out)+".bed", 9, 1000.0, MergeOutPrefix(PfstPrefix(set.Out)) + "_perc")
	if err != nil { return }
	err = ThreshAndMerge(FstPrefix(set.Out) + ".bed", 4, .10, MergeOutPrefix(FstPrefix(set.Out)) + "_perc")
	if err != nil { return }
	err = ThreshAndMerge(SelecPrefix(set.Out) + ".bed", 4, .04, MergeOutPrefix(SelecPrefix(set.Out)) + "_perc")
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

// func ParseGoodAndAlts(line []string) GoodAndAlts {
// 	var gaa GoodAndAlts
// 	if len(line) > 0 {
// 		gaa.Good = line[0]
// 	}
// 	gaa.Alts = slices.Clone(line[1:])
// 	return gaa
// }
// 
// func GetGoodAndAlts(r io.Reader) []GoodAndAlts {
// 	s := lscan.NewScanner(r, lscan.ByByte('\t'))
// 	var gaas []GoodAndAlts
// 	for s.Scan() {
// 		gaas = append(gaas, ParseGoodAndAlts(s.Line()))
// 	}
// 	return gaas
// }

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
}

func GetPlot4Flags() Plot4Flags {
	var f Plot4Flags

	flag.IntVar(&f.Threads, "t", 1, "Threads to use")
	percstring := flag.String("p", ".001", "Percentile threshold")
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

func SubtractAlts(gset GoodAndAlts, statistic string) (outpath string, err error) {
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
		good.SubtractBed(alt)
	}

	outpath = gset.Good.OutputPrefix(statistic) + ".bed"
	ovlconn, err := os.Create(outpath)
	if err != nil {
		return "", err
	}
	defer ovlconn.Close()

	permuvals.FprintBeds(ovlconn, good)

	return outpath, err
}

func SubtractAllAlts(gsets []GoodAndAlts, statistics []string, threads int) (outpaths []string, errs []error) {
	fmt.Println(statistics)
	var wg sync.WaitGroup
	njobs := len(gsets) * len(statistics)
	mod := len(statistics)
	jobs := make(chan int, njobs)
	outpaths = make([]string, njobs)
	errs = make([]error, njobs)

	for i:=0; i<threads; i++ {
		go func() {
			for job := range jobs {
				outpaths[job], errs[job] = SubtractAlts(gsets[job/mod], statistics[job%mod])
				wg.Done()
			}
		}()
	}

	wg.Add(njobs)
	for i:=0; i<njobs; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return outpaths, errs
}

func PercTMAParallel(set PlotSet, errchan chan error, percentile float64) {
	err := PercThreshMergeAll(set, percentile)
	errchan <- err
}
func PercTMASets(sets PlotSets, percentile float64) (errors []error) {
	var errorchans []chan error
	for _, set := range sets {
		ec := make(chan error)
		errorchans = append(errorchans, ec)
		go PercTMAParallel(set, ec, percentile)
	}
	for _, ec := range errorchans {
		errors = append(errors, <-ec)
	}
	return
}


func RunGood4Plots() {
	flags := GetPlot4Flags()
	plot_sets := ReadPlotSets(os.Stdin)
	runtime.GOMAXPROCS(flags.Threads)

	errors := PercTMASets(plot_sets, flags.Percentile)
	for _, err := range errors {
		if err != nil { panic(err) }
	}

	// all_bedpaths := GetAllBedpaths(plot_sets)
	// errors = PASetsLimited(all_bedpaths, statistics, flags.Threads)

	goodsAndAlts := PrebuiltGoodAndAlts()
	statistics := []string{"pFst", "Fst", "Selec"}
	outpaths, errors := SubtractAllAlts(goodsAndAlts, statistics, flags.Threads)
	for _, err := range errors {
		if err != nil { fmt.Fprintln(os.Stderr, err) }
	}

	pathsconn, err := os.Create("goodbedpaths.txt")
	if err != nil { panic(err) }
	defer pathsconn.Close()
	for _, path := range outpaths {
		fmt.Fprintln(pathsconn, path)
	}

	// goodsAndAlts, err := ReadPath(flags.GoodsAndAltsPath, GetGoodAndAlts)
	// if err != nil {
	// 	panic(err)
	// }

	// for _, stat := range statistics {
	// 	ovlpath := stat + "_overlaps.txt"
	// 	for _, gset := range goodsAndAlts {
	// 		err = SubtractAlts(gset, ovlpath)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 	}
	// }
}

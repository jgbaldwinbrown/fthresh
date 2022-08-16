package fthresh

import (
	"runtime"
	"flag"
	"fmt"
	"os"
	"github.com/jgbaldwinbrown/permuvals/pkg"
	"github.com/jgbaldwinbrown/pmap/pkg"
)

func ThreshMergeAll(set PlotSet) (err error) {
	err = ThreshAndMerge(PfstPrefix(set.Out)+".bed", 9, 1000.0, MergeOutPrefix(PfstPrefix(set.Out)))
	if err != nil { return }
	err = ThreshAndMerge(FstPrefix(set.Out) + ".bed", 4, .10, MergeOutPrefix(FstPrefix(set.Out)))
	if err != nil { return }
	err = ThreshAndMerge(SelecPrefix(set.Out) + ".bed", 4, .04, MergeOutPrefix(SelecPrefix(set.Out)))
	if err != nil { return }

	return nil
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
	flags.Verbose = false
	flags.MaxComps = 3
	fmt.Println("start comp")
	comp, err := permuvals.FullCompare(flags)
	fmt.Println("end comp")
	if err != nil { return }

	outconn, err := os.Create(statistic + "_probs.txt")
	if err != nil { return }
	defer outconn.Close()

	permuvals.FprintProbs(outconn, comp.Probs)

	ovlconn, err := os.Create(statistic + "_overlaps.txt")
	if err != nil { return }
	defer ovlconn.Close()

	permuvals.FprintOvlsBed(ovlconn, comp.Overlaps)

	ovlstatsconn, err := os.Create(statistic + "_overlap_stats.txt")
	if err != nil { return }
	defer ovlstatsconn.Close()

	permuvals.FprintOvlsStats(ovlstatsconn, permuvals.AllOvlsStats(comp.Overlaps)...)

	return nil
}

func PAParallel(bedpaths []string, statistic string, seed int, errchan chan error) {
	err := PermAll(bedpaths, statistic, seed)
	errchan <- err
}

func TMASets(sets PlotSets) (errors []error) {
	return pmap.Map(ThreshMergeAll, sets, -1)
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

type seededBedpaths struct {
	seed int
	bedpaths []string
}

func PASetsLimited(all_bedpaths [][]string, statistics []string, threads int) (errors []error) {
	njobs := len(all_bedpaths)
	jobq := make(chan seededBedpaths, njobs)
	errs := make(chan error, njobs)
	for i:=0; i<threads; i++ {
		go func() {
			for sbp := range jobq {
				PAParallel(sbp.bedpaths, statistics[sbp.seed], sbp.seed, errs)
			}
		}()
	}
	for seed, bedpaths := range all_bedpaths {
		jobq <- seededBedpaths{seed, bedpaths}
	}
	close(jobq)
	for i := 0; i<njobs; i++ {
		errors = append(errors, <-errs)
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

type RvrFlags struct {
	Threads int
}

func GetRvrFlags() RvrFlags {
	var f RvrFlags
	flag.IntVar(&f.Threads, "t", 1, "Threads to use")
	flag.Parse()
	return f
}

func RunThreshMergePermRvr() {
	plot_sets := ReadPlotSets(os.Stdin)
	flags := GetRvrFlags()
	runtime.GOMAXPROCS(flags.Threads)

	errors := TMASets(plot_sets)
	for _, err := range errors {
		if err != nil { panic(err) }
	}

	all_bedpaths := GetAllBedpaths(plot_sets)
	statistics := []string{"pFst", "Fst", "Selec"}
	errors = PASetsLimited(all_bedpaths, statistics, flags.Threads)
	for _, err := range errors {
		if err != nil { panic(err) }
	}
}

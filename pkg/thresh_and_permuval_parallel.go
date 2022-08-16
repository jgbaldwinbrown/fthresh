package fthresh

import (
	"fmt"
	"os"
)

func ThreshMergeAndPerm(set PlotSet, statistic string, use_set_prefix bool, seed int) (err error) {
	err = ThreshMergeAll(set)
	if err != nil { return }

	bedpaths := []string{PfstMergeOut(set.Out), FstMergeOut(set.Out), SelecMergeOut(set.Out)}
	if use_set_prefix {
		statistic = statistic + set.Out
	}
	err = PermAll(bedpaths, statistic, seed)
	if err != nil { return }

	return nil
}

func TMPParallel(set PlotSet, statistic string, use_set_prefix bool, seed int, errchan chan error) {
	err := ThreshMergeAndPerm(set, statistic, use_set_prefix, seed)
	errchan <- err
}

func TMPSets(sets PlotSets, statistic string, use_set_prefix bool) (errors []error) {
	var errorchans []chan error
	for seed, set := range sets {
		ec := make(chan error)
		errorchans = append(errorchans, ec)
		go TMPParallel(set, statistic, use_set_prefix, seed, ec)
	}
	for _, ec := range errorchans {
		errors = append(errors, <-ec)
	}
	return
}

func ThreshMergeAndPermParallel() {
	plot_sets := ReadPlotSets(os.Stdin)
	errors := TMPSets(plot_sets, "pfst_allcombo", true)
	for _, err := range errors {
		// if err != nil { panic(err) }
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

package fthresh

import (
	"os"
)

func ThreshMergeAndPerm(set PlotSet, statistic string, seed int) (err error) {
	err = ThreshMergeAll(set)
	if err != nil { return }

	bedpaths := []string{PfstMergeOut(set.Out), FstMergeOut(set.Out), SelecMergeOut(set.Out)}
	err = PermAll(bedpaths, statistic, seed)
	if err != nil { return }

	return nil
}

func TMPParallel(set PlotSet, statistic string, seed int, errchan chan error) {
	err := ThreshMergeAndPerm(set, statistic, seed)
	errchan <- err
}

func TMPSets(sets PlotSets, statistic string) (errors []error) {
	var errorchans []chan error
	for seed, set := range sets {
		ec := make(chan error)
		errorchans = append(errorchans, ec)
		go TMPParallel(set, statistic, seed, ec)
	}
	for _, ec := range errorchans {
		errors = append(errors, <-ec)
	}
	return
}

func ThreshMergeAndPermParallel() {
	plot_sets := ReadPlotSets(os.Stdin)
	errors := TMPSets(plot_sets, "nostat")
	for _, err := range errors {
		if err != nil { panic(err) }
	}
}

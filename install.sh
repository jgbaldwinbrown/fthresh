#!/bin/bash
set -e

go mod tidy

(cd cmd && (
	ls -1d *.go | while read i ; do
		go build $i
	done
))

cp cmd/thresh_and_permuval ~/mybin/
cp cmd/thresh_and_permuval_rvr ~/mybin/
cp cmd/combine_fst_and_pfsts ~/mybin/
cp cmd/make_pfst_fst_selec ~/mybin/
cp cmd/combine_pfst_fst_selec ~/mybin/
cp cmd/thresh_and_permuval_parallel ~/mybin/
cp cmd/qqplot_all ~/mybin/
cp cmd/good4plots ~/mybin/
cp cmd/plot_good4plots ~/mybin/
cp cmd/thresh_merge_and_gggenes ~/mybin/

(cd scripts && (
	ls -1d *.go | while read i ; do
		go build $i
	done
))

cp scripts/make_all_plots ~/mybin/
cp scripts/lfmt ~/mybin/
cp scripts/plfmt_flex ~/mybin/plfmt_flex
cp scripts/make_all_multiplots ~/mybin/

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
cp cmd/good2plots ~/mybin/
cp cmd/plot_good4plots ~/mybin/
cp cmd/plot_good4plots_together ~/mybin/
cp cmd/plot_good4plots_vert ~/mybin/
cp cmd/plot_good2plots ~/mybin/
cp cmd/plot_good2plots_together ~/mybin/
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
cp scripts/plot_pfst_fst_selec.R ~/mybin/plot_pfst_fst_selec
chmod +x ~/mybin/plot_pfst_fst_selec
cp scripts/plot_pfst_fst.R ~/mybin/plot_pfst_fst
chmod +x ~/mybin/plot_pfst_fst
cp scripts/plot_goods.R ~/mybin/plot_goods
chmod +x ~/mybin/plot_goods
cp scripts/plot_goods_pfst_together.R ~/mybin/plot_goods_pfst_together
chmod +x ~/mybin/plot_goods_pfst_together
cp scripts/plot_goods_pfst_together2.R ~/mybin/plot_goods_pfst_together2
chmod +x ~/mybin/plot_goods_pfst_together2
cp scripts/plot_goods_pfst_vert.R ~/mybin/plot_goods_pfst_vert
chmod +x ~/mybin/plot_goods_pfst_vert
cp scripts/plot_pretty_multiple_helpers.R ~/rlibs

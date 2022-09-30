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
cp cmd/make_pfst_fst_selec_nowin ~/mybin/
cp cmd/combine_pfst_fst_selec ~/mybin/
cp cmd/thresh_and_permuval_parallel ~/mybin/
cp cmd/qqplot_all ~/mybin/
cp cmd/good4plots ~/mybin/
cp cmd/good4plots_reps ~/mybin/
cp cmd/good2plots ~/mybin/
cp cmd/plot_good4plots ~/mybin/
cp cmd/plot_good4plots_reps ~/mybin/
cp cmd/plot_good4plots_together ~/mybin/
cp cmd/plot_good4plots_reps_together_reptop ~/mybin/
cp cmd/plot_good4plots_reps_together ~/mybin/
cp cmd/plot_good4plots_vert ~/mybin/
cp cmd/plot_good8plots_vert ~/mybin/
cp cmd/plot_good2plots ~/mybin/
cp cmd/plot_good2plots_together ~/mybin/
cp cmd/thresh_merge_and_gggenes ~/mybin/
# cp cmd/gggenes_noselec ~/mybin/

(cd scripts && (
	ls -1d *.go | while read i ; do
		go build $i
	done
))

cp scripts/make_all_plots ~/mybin/
cp scripts/make_all_plots_nowin ~/mybin/
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
cp scripts/plot_goods_pfst_together_nowin.R ~/mybin/plot_goods_pfst_together_nowin
chmod +x ~/mybin/plot_goods_pfst_together_nowin
cp scripts/plot_goods_pfst_reps_together.R ~/mybin/plot_goods_pfst_reps_together
chmod +x ~/mybin/plot_goods_pfst_reps_together
cp scripts/plot_goods_pfst_reps_together_reptop.R ~/mybin/plot_goods_pfst_reps_together_reptop
chmod +x ~/mybin/plot_goods_pfst_reps_together_reptop
cp scripts/plot_goods_pfst_reps_together_reptop_nowin.R ~/mybin/plot_goods_pfst_reps_together_reptop_nowin
chmod +x ~/mybin/plot_goods_pfst_reps_together_reptop_nowin
cp scripts/plot_goods_pfst_together2.R ~/mybin/plot_goods_pfst_together2
chmod +x ~/mybin/plot_goods_pfst_together2
cp scripts/plot_goods_pfst_vert.R ~/mybin/plot_goods_pfst_vert
chmod +x ~/mybin/plot_goods_pfst_vert
cp scripts/plot_goods_pfst_vert_size.R ~/mybin/plot_goods_pfst_vert_size
chmod +x ~/mybin/plot_goods_pfst_vert_size
cp scripts/plot_region.R ~/mybin/plot_region
chmod +x ~/mybin/plot_region
cp scripts/genelists.sh ~/mybin/genelists
chmod +x ~/mybin/genelists
cp scripts/bedify.sh ~/mybin/bedify
chmod +x ~/mybin/bedify
cp scripts/fdr_it.py ~/mybin/fdr_it
chmod +x ~/mybin/fdr_it
cp scripts/plfmt.py ~/mybin/plfmt
chmod +x ~/mybin/plfmt
cp scripts/all_winplot_multipheno.sh ~/mybin/all_winplot_multipheno
chmod +x ~/mybin/all_winplot_multipheno
cp scripts/prep_all_full.sh ~/mybin/prep_all_full
chmod +x ~/mybin/prep_all_full
cp scripts/good4plotsfullwrap.sh ~/mybin/good4plotsfullwrap
chmod +x ~/mybin/good4plotsfullwrap

cp scripts/plot_pretty_multiple_helpers.R ~/rlibs

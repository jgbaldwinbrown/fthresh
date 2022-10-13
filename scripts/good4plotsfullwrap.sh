#!/bin/bash
set -e

WINSIZE=50000
WINSTEP=5000
WINSTR=""

if [ $# -ge 2 ] ; then
	WINSIZE=${1}
	WINSTEP=${2}
	WINSTR=_win${WINSIZE}_${WINSTEP}
fi

(
	echo "bash prep_all_full.sh"
	if [ $# -ge 2 ] ; then
		prep_all_full $WINSIZE $WINSTEP
	else
		prep_all_full
	fi

	echo "bash pfst_fst_selec_list_full.sh"
	combine_pfst_fst_selec_json \
		all_bigfile_fsts${WINSTR}_full.txt \
		all_bigfile_pfsts${WINSTR}_full.txt \
		old_selecs_1k_paths.txt \
		"${WINSIZE},${WINSTEP}" \
	> pfst_fst_selec_combo${WINSTR}_full_json.txt

	echo "zero"

	cat pfst_fst_selec_combo${WINSTR}_full_json.txt | make_pfst_fst_selec louse_genome_0.1.1.chrlens.bed

	echo "one"

	cat pfst_fst_selec_combo${WINSTR}_full_json.txt | \
	good4plots -p 0.001 -t 4 -f

	cat pfst_fst_selec_combo${WINSTR}_full_json.txt | \
	good4plots_reps -p 0.001 -t 4 -f

	cat pfst_fst_selec_combo${WINSTR}_full_json.txt | \
	plot_good4plots louse_genome_0.1.1.chrlens.bed

	cat pfst_fst_selec_combo${WINSTR}_full_json.txt | \
	plot_good4plots_together -f louse_genome_0.1.1.chrlens.bed \
	> plot_good4plotsfull_script${WINSTR}.sh


	# cat pfst_fst_selec_combo${WINSTR}_full_json.txt | \
	# plot_good4plots_reps louse_genome_0.1.1.chrlens.bed

	echo "four"

	bash plot_good4plotsfull_script${WINSTR}.sh

	echo "four point one"

	cat pfst_fst_selec_combo${WINSTR}_full_json.txt | \
	plot_good4plotsfull_reps_together_reptop -f louse_genome_0.1.1.chrlens.bed \
	> plot_good4plotsfull_reps_reptop_script${WINSTR}.sh

	echo "four point two"

	bash plot_good4plots_reps_reptop_script${WINSTR}.sh

	ls *subtractedalt*subfull*plfmt*bed | genelists louse_annotation_0.1.1.gff
	ls *subtractedalt*subfull*plfmt*bed*genes*gff | grep 'fullrep' | xargs cat > all_genes_fullrep.gff
	cat all_genes_fullrep.gff | grep -o 'Note[^;]*;' | sed 's/Note=Similar to \|;//g' > all_gene_descriptions_fullrep.txt

	# echo "eight"

	cat pfst_fst_selec_combo${WINSTR}_full_json.txt | \
	gggenes_noselec \
	> gggenes_noselec_out_fullrep${WINSTR}.txt
) > good4plotsfull_out${WINSTR}.txt 2>&1

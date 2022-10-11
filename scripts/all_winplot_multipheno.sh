#!/bin/bash
set -e

WINSIZE=50000
WINSTEP=5000
WINSTR="_win"

if [ $# -ge 2 ] ; then
	WINSIZE="${1}"
	WINSTEP="${2}"
	WINSTR="_win${WINSIZE}_${WINSTEP}"
fi

find allnames/ -name '*.fst' | sort | while read i ; do
	if [ -s "${i}" ] ; then
		# a=`basename "${i}"`
		# echo "cat ${i} | fst_sliding_window -w 50000 -s 5000 > ${i}_win.txt && \
		# python3 plot_vcftools_fst_win.py ${i}_win.txt louse_genome_0.1.1.chrlens.bed ${i}_win_plot ${a}_win && \
		# convert -density 100 ${i}_win_plot.pdf -compress lzw ${i}_win_plot.png"
		echo "cat ${i} | fst_sliding_window -w ${WINSIZE} -s ${WINSTEP} > ${i}${WINSTR}.txt"
	fi
done | \
parallel -j 8 {}

# ls all_comparisons_out/*fst | while read i ; do
# 	a=`basename "${i}"`
# 	echo "python3 plot_vcftools_fst.py ${i} louse_genome_0.1.1.chrlens.bed all_comparisons_out/${a}_plot ${a}"
# done | \
#  parallel -j 8 {}

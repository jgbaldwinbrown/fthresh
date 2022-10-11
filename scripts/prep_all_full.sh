#!/bin/bash
set -e

#BIGFILEPATH=/media/jgbaldwinbrown/3564-3063/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns
BIGFILEPATH=/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns
CWD=`pwd`
WINSIZE=50000
WINSTEP=5000
WINSTR=""

if [ $# -ge 2 ] ; then
	WINSIZE=$1
	WINSTEP=$2
	WINSTR=_win${WINSIZE}_${WINSTEP}
fi
FSTPATH="${CWD}/all_bigfile_fsts${WINSTR}_full.txt"
PFSTPATH="${CWD}/all_bigfile_pfsts${WINSTR}_full.txt"

cp window_fisher_bp.py "${BIGFILEPATH}"
cp bedify.sh "${BIGFILEPATH}"
cp fdr.py "${BIGFILEPATH}"
cp fdr_it.py "${BIGFILEPATH}"
cp plfmt.py "${BIGFILEPATH}"
cp plot_pretty_hlines_bp.R "${BIGFILEPATH}"

cd "${BIGFILEPATH}" && (
	echo "one"
	if [ $# -ge 2 ] ; then
		all_winplot_multipheno $WINSIZE $WINSTEP
	else
		all_winplot_multipheno
	fi
	echo "two"

	if [ $# -ge 2 ] ; then
		find . -name '*pfst.txt' | \
		grep -v 'backup' | \
		grep -v 'Makefile' | \
		sed 's/\.txt$//' | \
		make_all_plots -w $WINSIZE -s $WINSTEP > Makefile_prep_pfst${WINSTR}
	else
		find . -name '*pfst.txt' | \
		grep -v 'backup' | \
		grep -v 'Makefile' | \
		sed 's/\.txt$//' | \
		make_all_plots > Makefile_prep_pfst${WINSTR}
	fi

	make -j 8 -k -f Makefile_prep_pfst${WINSTR} > make_out_Makefile_prep_pfst${WINSTR}.txt 2>&1

	echo "three"

	if [ $# -ge 2 ] ; then
		find $PWD -name '*\.fst'"${WINSTR}"'.txt' | grep -v 'backup' > "${FSTPATH}"
		echo "four"
		find $PWD -name '*pfst'"${WINSTR}"'_fdr.bed' | grep -v 'backup' > "${PFSTPATH}"
		echo "five"
	else
		find $PWD -name '*\.fst_win.txt' | grep -v 'backup' > "${FSTPATH}"
		echo "four"
		find $PWD -name '*pfst_win_fdr.bed' | grep -v 'backup' > "${PFSTPATH}"
		echo "five"
	fi
)

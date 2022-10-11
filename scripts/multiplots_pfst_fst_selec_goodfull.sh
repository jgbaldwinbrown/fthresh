#!/bin/bash
set -e

# combine_pfst_fst_selec \
# 	all_bigfile_fsts.txt \
# 	all_bigfile_pfsts.txt \
# 	old_selecs_1k_paths.txt \
# > pfst_fst_selec_combo.txt

cat goodfull_pfst_fst_selec_combo.txt | make_pfst_fst_selec louse_genome_0.1.1.chrlens.bed
cat goodfull_pfst_fst_selec_combo_rep.txt | make_pfst_fst_selec louse_genome_0.1.1.chrlens.bed

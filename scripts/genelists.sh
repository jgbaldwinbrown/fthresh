#!/bin/bash
set -e

while read i ; do
	bedtools intersect -a "${1}" -b "${i}" > "${i}_features.gff"
	cat "${i}_features.gff" | mawk -F "\t" -v OFS="\t" '$3=="gene"' > ${i}_genes.gff
done

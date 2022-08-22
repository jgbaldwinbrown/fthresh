#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_pretty_multiple_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)
	print("args:")
	print(args)
	pfst_path = args[1]
	fst_path = args[2]
	selec_path = args[3]

	pfst_sig_path = args[4]
	# fst_sig_path = args[6]
	# selec_sig_path = args[7]

	out_path = args[5]

	pfst = read_combined_pvals_precomputed(pfst_path)
	fst = read_bed(fst_path, FALSE)
	selec = read_selec(selec_path)

	pfst_hithresh = calc_thresh(pfst, "VAL", .9999, TRUE)
	pfst_lowthresh = -log10(0.05)

	pfst$pass = pass_thresh(pfst, "VAL", pfst_hithresh)
	pfst$color = threshcolor(pfst, "CHR", "pass")

	pfst_threshes = data.frame(THRESH=c(pfst_hithresh, pfst_lowthresh))

	fst$color = nothreshcolor(fst, "CHR")

	selec$color = nothreshcolor(selec, "CHR")

	joinlist = join(list(pfst, fst, selec), list(pfst_threshes), c("pfst", "fst", "selec"), "pfst")
	data = joinlist[[1]]
	thresholds = joinlist[[2]]

	pfst_sig_rect = bed2rect(pfst_sig_path)
	pfst_sig_rect$NAME = "pfst"
	# fst_sig_rect = bed2rect(fst_sig_path)
	# selec_sig_rect = bed2rect(selec_sig_path)

	significant_boxes = pfst_sig_rect

	scales_y = list (
		`pfst` = scale_y_continuous(limits = c(0, 350)),
		`fst` = scale_y_continuous(limits = c(0, .3)),
		`selec` = scale_y_continuous(limits = c(-.01, .06))
	)

	print("plotting to:")
	print(out_path)

	print("rects:")
	print(significant_boxes)
	plot_scaled_y_boxed(data, VAL, out_path, 20, 8, 300, thresholds, calc_chrom_labels(pfst), scales_y, significant_boxes)
}

main()

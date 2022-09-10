#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_pretty_multiple_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	figurita_pfst_path = args[1]
	figurita_pfst_sig_path = args[2]

	runt_pfst_path = args[3]
	runt_pfst_sig_path = args[4]

	figurita_bitted_pfst_path = args[5]
	figurita_bitted_pfst_sig_path = args[6]

	runt_bitted_pfst_path = args[7]
	runt_bitted_pfst_sig_path = args[8]

	out_path = args[9]

	figurita_pfst = read_combined_pvals_precomputed(figurita_pfst_path)
	runt_pfst = read_combined_pvals_precomputed(runt_pfst_path)
	figurita_bitted_pfst = read_combined_pvals_precomputed(figurita_bitted_pfst_path)
	runt_bitted_pfst = read_combined_pvals_precomputed(runt_bitted_pfst_path)

	figurita_pfst_hithresh = calc_thresh(figurita_pfst, "VAL", .999, TRUE)
	figurita_pfst_lowthresh = -log10(0.05)
	figurita_pfst$pass = pass_thresh(figurita_pfst, "VAL", figurita_pfst_hithresh)
	figurita_pfst$color = threshcolor(figurita_pfst, "CHR", "pass")
	figurita_pfst_threshes = data.frame(THRESH=c(figurita_pfst_hithresh, figurita_pfst_lowthresh))

	runt_pfst_hithresh = calc_thresh(runt_pfst, "VAL", .999, TRUE)
	runt_pfst_lowthresh = -log10(0.05)
	runt_pfst$pass = pass_thresh(runt_pfst, "VAL", runt_pfst_hithresh)
	runt_pfst$color = threshcolor(runt_pfst, "CHR", "pass")
	runt_pfst_threshes = data.frame(THRESH=c(runt_pfst_hithresh, runt_pfst_lowthresh))

	figurita_bitted_pfst_hithresh = calc_thresh(figurita_bitted_pfst, "VAL", .999, TRUE)
	figurita_bitted_pfst_lowthresh = -log10(0.05)
	figurita_bitted_pfst$pass = pass_thresh(figurita_bitted_pfst, "VAL", figurita_bitted_pfst_hithresh)
	figurita_bitted_pfst$color = threshcolor(figurita_bitted_pfst, "CHR", "pass")
	figurita_bitted_pfst_threshes = data.frame(THRESH=c(figurita_bitted_pfst_hithresh, figurita_bitted_pfst_lowthresh))

	runt_bitted_pfst_hithresh = calc_thresh(runt_bitted_pfst, "VAL", .999, TRUE)
	runt_bitted_pfst_lowthresh = -log10(0.05)
	runt_bitted_pfst$pass = pass_thresh(runt_bitted_pfst, "VAL", runt_bitted_pfst_hithresh)
	runt_bitted_pfst$color = threshcolor(runt_bitted_pfst, "CHR", "pass")
	runt_bitted_pfst_threshes = data.frame(THRESH=c(runt_bitted_pfst_hithresh, runt_bitted_pfst_lowthresh))

	joinlist = join(
		list(figurita_pfst, runt_pfst, figurita_bitted_pfst, runt_bitted_pfst),
		list(figurita_pfst_threshes, runt_pfst_threshes, figurita_bitted_pfst_threshes, runt_bitted_pfst_threshes),
		c("figurita", "runt", "figurita_bitted", "runt_bitted"),
		c("figurita", "runt", "figurita_bitted", "runt_bitted")
	)
	data = joinlist[[1]]
	thresholds = joinlist[[2]]

	scales_y = list (
		`figurita` = scale_y_continuous(limits = c(0, 350)),
		`runt` = scale_y_continuous(limits = c(0, 350)),
		`figurita_bitted` = scale_y_continuous(limits = c(0, 350)),
		`runt_bitted` = scale_y_continuous(limits = c(0, 350))
	)

	plot_scaled_y_vert(data, VAL, out_path, 20, 8, 300, thresholds, calc_chrom_labels(figurita_pfst), scales_y)
}

main()

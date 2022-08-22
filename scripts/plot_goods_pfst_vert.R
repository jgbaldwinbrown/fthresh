#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_pretty_multiple_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)

	black_pfst_path = args[1]
	black_pfst_sig_path = args[2]

	white_pfst_path = args[3]
	white_pfst_sig_path = args[4]

	black_bitted_pfst_path = args[5]
	black_bitted_pfst_sig_path = args[6]

	white_bitted_pfst_path = args[7]
	white_bitted_pfst_sig_path = args[8]

	out_path = args[9]

	black_pfst = read_combined_pvals_precomputed(black_pfst_path)
	white_pfst = read_combined_pvals_precomputed(white_pfst_path)
	black_bitted_pfst = read_combined_pvals_precomputed(black_bitted_pfst_path)
	white_bitted_pfst = read_combined_pvals_precomputed(white_bitted_pfst_path)

	black_pfst_hithresh = calc_thresh(black_pfst, "VAL", .999, TRUE)
	black_pfst_lowthresh = -log10(0.05)
	black_pfst$pass = pass_thresh(black_pfst, "VAL", black_pfst_hithresh)
	black_pfst$color = threshcolor(black_pfst, "CHR", "pass")
	black_pfst_threshes = data.frame(THRESH=c(black_pfst_hithresh, black_pfst_lowthresh))

	white_pfst_hithresh = calc_thresh(white_pfst, "VAL", .999, TRUE)
	white_pfst_lowthresh = -log10(0.05)
	white_pfst$pass = pass_thresh(white_pfst, "VAL", white_pfst_hithresh)
	white_pfst$color = threshcolor(white_pfst, "CHR", "pass")
	white_pfst_threshes = data.frame(THRESH=c(white_pfst_hithresh, white_pfst_lowthresh))

	black_bitted_pfst_hithresh = calc_thresh(black_bitted_pfst, "VAL", .999, TRUE)
	black_bitted_pfst_lowthresh = -log10(0.05)
	black_bitted_pfst$pass = pass_thresh(black_bitted_pfst, "VAL", black_bitted_pfst_hithresh)
	black_bitted_pfst$color = threshcolor(black_bitted_pfst, "CHR", "pass")
	black_bitted_pfst_threshes = data.frame(THRESH=c(black_bitted_pfst_hithresh, black_bitted_pfst_lowthresh))

	white_bitted_pfst_hithresh = calc_thresh(white_bitted_pfst, "VAL", .999, TRUE)
	white_bitted_pfst_lowthresh = -log10(0.05)
	white_bitted_pfst$pass = pass_thresh(white_bitted_pfst, "VAL", white_bitted_pfst_hithresh)
	white_bitted_pfst$color = threshcolor(white_bitted_pfst, "CHR", "pass")
	white_bitted_pfst_threshes = data.frame(THRESH=c(white_bitted_pfst_hithresh, white_bitted_pfst_lowthresh))

	joinlist = join(
		list(black_pfst, white_pfst, black_bitted_pfst, white_bitted_pfst),
		list(black_pfst_threshes, white_pfst_threshes, black_bitted_pfst_threshes, white_bitted_pfst_threshes),
		c("black", "white", "black_bitted", "white_bitted"),
		c("black", "white", "black_bitted", "white_bitted")
	)
	data = joinlist[[1]]
	thresholds = joinlist[[2]]

	scales_y = list (
		`black` = scale_y_continuous(limits = c(0, 350)),
		`white` = scale_y_continuous(limits = c(0, 350)),
		`black_bitted` = scale_y_continuous(limits = c(0, 350)),
		`white_bitted` = scale_y_continuous(limits = c(0, 350))
	)

	plot_scaled_y_vert(data, VAL, out_path, 20, 8, 300, thresholds, calc_chrom_labels(black_pfst), scales_y)
}

main()

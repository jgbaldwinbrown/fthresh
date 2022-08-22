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

	out_path = args[5]

	black_pfst = read_combined_pvals_precomputed(black_pfst_path)
	white_pfst = read_combined_pvals_precomputed(white_pfst_path)

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

	joinlist = join(
		list(black_pfst, white_pfst),
		list(black_pfst_threshes, white_pfst_threshes),
		c("black", "white"),
		c("black", "white")
	)
	data = joinlist[[1]]
	thresholds = joinlist[[2]]

	black_pfst_sig_rect = bed2rect(black_pfst_sig_path)
	black_pfst_sig_rect$NAME = "black"
	white_pfst_sig_rect = bed2rect(white_pfst_sig_path)
	white_pfst_sig_rect$NAME = "white"

	significant_boxes = as.data.frame(rbind(
		black_pfst_sig_rect,
		white_pfst_sig_rect
	))

	scales_y = list (
		`black` = scale_y_continuous(limits = c(0, 350)),
		`white` = scale_y_continuous(limits = c(0, 350))
	)

	print("boxes:")
	print(significant_boxes)
	plot_scaled_y_boxed(data, VAL, out_path, 20, 8, 300, thresholds, calc_chrom_labels(black_pfst), scales_y, significant_boxes)
}

main()

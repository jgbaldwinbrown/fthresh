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

	runt_pfst_path = args[5]
	runt_pfst_sig_path = args[6]

	figurita_pfst_path = args[7]
	figurita_pfst_sig_path = args[8]

	out_path = args[9]

	black_pfst = read_combined_pvals_precomputed(black_pfst_path)
	white_pfst = read_combined_pvals_precomputed(white_pfst_path)
	runt_pfst = read_combined_pvals_precomputed(runt_pfst_path)
	figurita_pfst = read_combined_pvals_precomputed(figurita_pfst_path)

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

	runt_pfst_hithresh = calc_thresh(runt_pfst, "VAL", .999, TRUE)
	runt_pfst_lowthresh = -log10(0.05)
	runt_pfst$pass = pass_thresh(runt_pfst, "VAL", runt_pfst_hithresh)
	runt_pfst$color = threshcolor(runt_pfst, "CHR", "pass")
	runt_pfst_threshes = data.frame(THRESH=c(runt_pfst_hithresh, runt_pfst_lowthresh))

	figurita_pfst_hithresh = calc_thresh(figurita_pfst, "VAL", .999, TRUE)
	figurita_pfst_lowthresh = -log10(0.05)
	figurita_pfst$pass = pass_thresh(figurita_pfst, "VAL", figurita_pfst_hithresh)
	figurita_pfst$color = threshcolor(figurita_pfst, "CHR", "pass")
	figurita_pfst_threshes = data.frame(THRESH=c(figurita_pfst_hithresh, figurita_pfst_lowthresh))

	joinlist = join(
		list(black_pfst, white_pfst, runt_pfst, figurita_pfst),
		list(black_pfst_threshes, white_pfst_threshes, runt_pfst_threshes, figurita_pfst_threshes),
		c("black", "white", "runt", "figurita"),
		c("black", "white", "runt", "figurita")
	)
	data = joinlist[[1]]
	thresholds = joinlist[[2]]

	black_pfst_sig_rect = bed2rect(black_pfst_sig_path)
	if (nrow(black_pfst_sig_rect) > 0) {
		black_pfst_sig_rect$NAME = "black"
	}
	white_pfst_sig_rect = bed2rect(white_pfst_sig_path)
	if (nrow(white_pfst_sig_rect) > 0) {
		white_pfst_sig_rect$NAME = "white"
	}
	runt_pfst_sig_rect = bed2rect(runt_pfst_sig_path)
	if (nrow(runt_pfst_sig_rect) > 0) {
		runt_pfst_sig_rect$NAME = "runt"
	}
	figurita_pfst_sig_rect = bed2rect(figurita_pfst_sig_path)
	if (nrow(figurita_pfst_sig_rect) > 0) {
		figurita_pfst_sig_rect$NAME = "figurita"
	}

	print("pre sig boxes")

	significant_boxes = as.data.frame(rbind(
		black_pfst_sig_rect,
		white_pfst_sig_rect,
		runt_pfst_sig_rect,
		figurita_pfst_sig_rect
	))

	print("pre scales_y")

	scales_y = list (
		`black` = scale_y_continuous(limits = c(0, 350)),
		`white` = scale_y_continuous(limits = c(0, 350)),
		`runt` = scale_y_continuous(limits = c(0, 350)),
		`figurita` = scale_y_continuous(limits = c(0, 350))
	)

	print("pre text")

	text = data.frame(
		x = c(0,0,0,0),
		y = c(350,350,350,350),
		NAME = c("black", "white", "runt", "figurita"),
		textlabel = c("Color (dark)", "Color (light)", "Size (large)", "Size (small)")
	)

	print("pre plot")

	plot_scaled_y_boxed_text(data, VAL, out_path, 20, 8, 300, thresholds, calc_chrom_labels(black_pfst), scales_y, significant_boxes, text)
}

main()

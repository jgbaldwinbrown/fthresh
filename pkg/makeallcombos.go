package fthresh
// 
// import (
// 	"fmt"
// 	"regexp"
// 	"os"
// 	"io"
// 	"bufio"
// 	"strings"
// 	"flag"
// )
// 
// func ReadLines(path string) ([]string, error) {
// 	r, err := os.Open(path)
// 	if err != nil {
// 		return err
// 	}
// 	var lines []string
// 	s := bufio.NewScanner(r)
// 	s.Buffer([]byte{}, 1e12)
// 	for s.Scan() {
// 		lines = append(lines(s.Text()))
// 	}
// 	return lines, nil
// }
// 
// type Maker func(lines []string, winstr string) error
// 
// func MakeCombo(lines []string, winstr string) error {
// 	outpath := fmt.Sprintln("good4plotsfull_combo%d.txt", winstr)
// 	var cfgs []ComboConfig
// 	for _, t := range Treatments() {
// 		re := t.Regex()
// 		lines := FindLines(lines, re)
// 		if len(lines) != 1 {
// 			return fmt.Errorf("regex for treatment %v found %v lines: %v", t, len(lines), lines)
// 		}
// 		cfgs := append(cfgs, MakeCfg(t, lines[0]))
// 	}
// 	return nil
// }
// 
// func main() {
// 	var winsize, winstep int
// 	flag.IntVar(&winsize, "w", 50000, "winsize")
// 	flag.IntVar(&winstep, "s", 5000, "winstep")
// 
// 	winstr := ""
// 	if winsize != 50000 && winstep != 5000 {
// 		winstr = fmt.Sprintf("win%d_%d", winsize, winstep)
// 	}
// 
// 	lines, err := ReadLines(fmt.Sprintf("pfst_fst_selec_combo%s_full.txt", winstr))
// 	if err != nil {
// 		panic(err)
// 	}
// 
// 	for _, f := range []Maker{MakeCombo, MakeComboControl, MakeRepComboControl, MakeRepCombo, MakeRepComboControlReptop} {
// 		err := f(lines, winstr)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 
// 	// # cat good4plotsfull_combo${WINSTR}.txt | \
// 	// # cat good4plotsfull_combo${WINSTR}_control.txt | \
// 	// # cat good4plotsfull_reps_combo${WINSTR}_control.txt | \
// 	// # cat good4plotsfull_reps_combo${WINSTR}.txt | \
// 	// # cat good4plotsfull_reps_combo${WINSTR}_control_reptop.txt | \
// }
// 
// 
// 
// /*
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Full_Full__pfst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Full_Full__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Full_Full__fst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Full_Full__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup3/black_pooled_unbitted_tle30_s_coeff_win1k.txt	_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Full_Full___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Low_Mid__pfst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Low_Mid__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Low_Mid__fst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Low_Mid__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup3/black_pooled_unbitted_tle30_s_coeff_win1k.txt	_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Low_Mid___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Full_Full__pfst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Full_Full__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Full_Full__fst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Full_Full__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup3/black_pooled_unbitted_tle30_s_coeff_win1k.txt	_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Full_Full___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Low_Mid__pfst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Low_Mid__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Low_Mid__fst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Low_Mid__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup3/black_pooled_unbitted_tle30_s_coeff_win1k.txt	_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Low_Mid___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Full_Full__pfst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Full_Full__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Full_Full__fst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Full_Full__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup3/black_pooled_unbitted_tle30_s_coeff_win1k.txt	_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Full_Full___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Low_High__pfst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Low_High__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Low_High__fst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Low_High__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup3/black_pooled_unbitted_tle30_s_coeff_win1k.txt	_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_Color_Low_High___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Full_Full__pfst/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Full_Full__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Full_Full__fst/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Full_Full__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup2/figurita_pooled_bitted_tle30_s_coeff_win1k.txt	_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Full_Full___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Low_Mid__pfst/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Low_Mid__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Low_Mid__fst/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Low_Mid__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup2/figurita_pooled_bitted_tle30_s_coeff_win1k.txt	_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Feral_time_36_bit_Bitted_replicate_All_Size_Low_Mid___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Full_Full__pfst/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Full_Full__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Full_Full__fst/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Full_Full__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup2/figurita_pooled_bitted_tle30_s_coeff_win1k.txt	_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Full_Full___multiplot
// /media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Low_High__pfst/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Low_High__pfst_win_fdr.bed	/media/jgbaldwinbrown/jim_work1/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Low_High__fst/_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Low_High__fst.weir.fst_win.txt	/home/jgbaldwinbrown/Documents/work_stuff/louse/s_estimation/partials/window/backup2/figurita_pooled_bitted_tle30_s_coeff_win1k.txt	_breed_Figurita_time_36_bit_Bitted_replicate_All_breed_Runt_time_36_bit_Bitted_replicate_All_Size_Low_High___multiplot
// 
// 	> pfst_fst_selec_combo${WINSTR}_full.txt
// 	# cat goodfull_pfst_fst_selec_combo${WINSTR}.txt | make_pfst_fst_selec louse_genome_0.1.1.chrlens.bed
// 	# cat goodfull_pfst_fst_selec_combo${WINSTR}_rep.txt | make_pfst_fst_selec louse_genome_0.1.1.chrlens.bed
// 	# cat good4plotsfull_combo${WINSTR}_control.txt | \
// 	# cat good4plotsfull_combo${WINSTR}.txt | \
// 	# cat good4plotsfull_combo${WINSTR}_control.txt | \
// 	# cat good4plotsfull_reps_combo${WINSTR}_control.txt | \
// 	# cat good4plotsfull_reps_combo${WINSTR}.txt | \
// 	# cat good4plotsfull_reps_combo${WINSTR}_control_reptop.txt | \
// 	# cat good4plotsfull_combo${WINSTR}_control.txt | \
// 
// 
// make_pfst_fst_selec
// good4plots
// plot_good4plots
// plot_good4plots_together
// good4plots_reps
// plot_good4plots_reps
// plot_good4plotsfull_reps_together_reptop
// gggenes_noselec
// 
// */

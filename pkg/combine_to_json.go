package fthresh

import (
	"regexp"
	"encoding/json"
	"fmt"
	"os"
	"io"
)

func (c Combo) ToConfig(winsize, winstep int) ComboConfig {
	var out ComboConfig
	ident := c.Fst.FstIdent()
	out.Treatment.Breed = ident.Breed
	out.Treatment.Bit = ident.Bit
	out.Treatment.Replicate = ident.Repl
	out.Treatment.Time = "36"

	out.WinSize = fmt.Sprintf("%d", winsize)
	out.WinStep = fmt.Sprintf("%d", winstep)
	winstr := fmt.Sprintf("win%d_%d", winsize, winstep)

	out.Fst = c.Fst.Path
	out.Pfst = c.Pfst.Path
	out.Selec = c.Selec.Path
	out.OutPrefix = c.Fst.Ident + winstr + "_multiplot"
	out.Subtractions = c.Fst.Ident + winstr + "_multiplot_subtractionsbed.bed"

	idents := c.Fst.FstIdents()
	if len(idents) > 1 {
		out.Treatment2.Breed = idents[1].Breed
		out.Treatment2.Bit = idents[1].Bit
		out.Treatment2.Replicate = idents[1].Repl
		out.Treatment2.Time = "36"

		out.ComparisonType = "nonsense"
		if out.Treatment.Breed != "feral" &&
			out.Treatment2.Breed == "feral" &&
			out.Treatment.Bit == out.Treatment2.Bit &&
			out.Treatment.Replicate == out.Treatment2.Replicate &&
			out.Treatment.Time == out.Treatment2.Time {
			if out.Treatment.Bit == "bitted" {
				out.ComparisonType = "control"
			} else {
				out.ComparisonType = "experimental"
			}
		}
	}

	return out
}

func ToConfig(winsize, winstep int, cs ...Combo) []ComboConfig {
	var out []ComboConfig
	fre := regexp.MustCompile("Full")
	for _, c := range cs {
		cfg := c.ToConfig(winsize, winstep)
		if fre.MatchString(cfg.Pfst) {
			out = append(out, c.ToConfig(winsize, winstep))
		}
	}
	return out
}

func FprintComboFstPfstJson(w io.Writer, combos []Combo, winsize, winstep int) error {
	cfgs := ToConfig(winsize, winstep, combos...)
	json, err := json.Marshal(cfgs)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n", json)
	return nil
}

func GetWinParams(arg string) (int, int, error) {
	var winsize, winstep int
	_, err := fmt.Sscanf(arg, "%d,%d", &winsize, &winstep)
	return winsize, winstep, err
}

func CombineFstPfstSelecJson() {
	all_fsts, err := GetPathFsts(os.Args[1])
	if err != nil {panic(err)}
	all_pfsts, err := GetPathPfsts(os.Args[2])
	if err != nil {panic(err)}
	all_selecs, err := GetPathSelecs(os.Args[3])
	if err != nil {panic(err)}

	winsize, winstep, err := GetWinParams(os.Args[4])
	if err != nil {panic(err)}

	combo := Combine(all_fsts, all_pfsts, all_selecs)
	err = FprintComboFstPfstJson(os.Stdout, combo, winsize, winstep)
	if err != nil {
		panic(err)
	}
}

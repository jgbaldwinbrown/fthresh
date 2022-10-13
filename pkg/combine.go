package fthresh

import (
	"fmt"
	"regexp"
	"os"
	"io"
	"bufio"
	"strings"
)

type Fst struct {
	Ident string
	Path string
}

type Combo struct {
	Fst Fst
	Pfst Fst
	Selec Fst
}

type Ident struct {
	Breed string
	Bit string
	Repl string
}

func GetFstsGeneric(fstre *regexp.Regexp, r io.Reader) (fs []Fst, err error) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)
	for s.Scan() {
		var f Fst
		f.Ident = fstre.FindString(s.Text())
		f.Path = s.Text()
		fs = append(fs, f)
	}
	return
}

func GetFsts(r io.Reader) ([]Fst, error) {
	fstre := regexp.MustCompile(`_breed_[^/]*__`)
	return GetFstsGeneric(fstre, r)
}

func GetPfsts(r io.Reader) ([]Fst, error) {
	return GetFsts(r)
}

func GetPathFsts(path string) ([]Fst, error) {
	fstsfile, err := os.Open(path)
	if err != nil { panic(err) }
	defer fstsfile.Close()
	return GetFsts(fstsfile)
}

func GetPathPfsts(path string) ([]Fst, error) {
	fstsfile, err := os.Open(path)
	if err != nil { panic(err) }
	defer fstsfile.Close()
	return GetPfsts(fstsfile)
}

func GetSelecs(r io.Reader) ([]Fst, error) {
	fstre := regexp.MustCompile(`[^/]*$`)
	return GetFstsGeneric(fstre, r)
}

func GetPathSelecs(path string) ([]Fst, error) {
	fstsfile, err := os.Open(path)
	if err != nil { panic(err) }
	defer fstsfile.Close()
	return GetSelecs(fstsfile)
}

func CombineFsts(fsts []Fst, pfsts []Fst) (combos []Combo) {
	for _, fst := range fsts {
		for _, pfst := range pfsts {
			if fst.Ident == pfst.Ident {
				combos = append(combos, Combo{fst, pfst, Fst{"", ""}})
			}
		}
	}
	return
}

func FprintComboFstPfst(w io.Writer, combos []Combo) {
	for _, c := range combos {
		fmt.Fprintf(w, "%v\t%v\t%v\n", c.Pfst.Path, c.Fst.Path, c.Fst.Ident + "_multiplot")
	}
}

func CombineFstPfst() {
	all_fsts, err := GetPathFsts(os.Args[1])
	if err != nil {panic(err)}
	all_pfsts, err := GetPathPfsts(os.Args[2])
	if err != nil {panic(err)}

	combo := CombineFsts(all_fsts, all_pfsts)
	FprintComboFstPfst(os.Stdout, combo)
}

////////////////////////////////////////////////////////////////////////////////////////////////////


func (i Ident) Match(i2 Ident) bool {
	return i.Breed == i2.Breed && i.Bit == i2.Bit && i.Repl == i2.Repl
}

func (i Ident) MatchExceptRepl(i2 Ident) bool {
	return i.Breed == i2.Breed && i.Bit == i2.Bit
}

func (f Fst) FstIdent() (out Ident) {
	re := regexp.MustCompile(`_breed_([a-zA-Z]*)_time_[^_]*_bit_([a-zA-Z]*)_replicate_([a-zA-Z0-9]*)`)
	matches := re.FindStringSubmatch(f.Ident)
	out.Breed = strings.ReplaceAll(strings.ToLower(matches[1]), "homer", "")
	out.Bit = strings.ToLower(matches[2])
	if matches[3] == "All" {
		out.Repl = matches[3]
	} else {
		out.Repl = strings.ReplaceAll(matches[3], "R", "")
	}
	return out
}
// /media/jgbaldwinbrown/3564-3063/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/from_laptop/vcftools_reruns/allnames/_breed_BlackHomer_time_36_bit_Unbitted_replicate_R1_breed_WhiteHomer_time_36_bit_Unbitted_replicate_R1_Color_Low_High__fst/_breed_BlackHomer_time_36_bit_Unbitted_replicate_R1_breed_WhiteHomer_time_36_bit_Unbitted_replicate_R1_Color_Low_High__fst.weir.fst_win.txt

func (f Fst) FstIdents() []Ident {
	re := regexp.MustCompile(`breed_([a-zA-Z]*)_time_[^_]*_bit_([a-zA-Z]*)_replicate_([a-zA-Z0-9]*)`)
	matches := re.FindAllStringSubmatch(f.Ident, -1)
	var outs []Ident
	for _, match := range matches {
		var out Ident
		out.Breed = strings.ReplaceAll(strings.ToLower(match[1]), "homer", "")
		out.Bit = strings.ToLower(match[2])
		if match[3] == "All" {
			out.Repl = match[3]
		} else {
			out.Repl = strings.ReplaceAll(match[3], "R", "")
		}
		outs = append(outs, out)
	}
	return outs
}

func (f Fst) SelecIdent() (out Ident) {
	re := regexp.MustCompile(`([a-z]*)_[a-z]*_([a-z]*).*(_repl([0-9]*))?`)
	matches := re.FindStringSubmatch(f.Ident)
	if len(matches) < 5 {
		return out
	}
	out.Breed = matches[1]
	out.Bit = matches[2]
	if matches[4] == "" {
		out.Repl = "All"
	} else {
		out.Repl = matches[4]
	}
	return out
}
// rk_stuff/louse/s_estimation/partials/window/white_pooled_unbitted_tle30_s_coeff_win1k.txt


func (f Fst) MatchSelec(selec Fst) bool {
	return f.FstIdent().MatchExceptRepl(selec.SelecIdent())
}

// rk_stuff/louse/s_estimation/partials/window/white_pooled_unbitted_tle30_s_coeff_win1k.txt
// func Combine(fsts []Fst, pfsts []Fst, selecs []Fst) (combos []Combo) {
// 	for _, fst := range fsts {
// 		for _, pfst := range pfsts {
// 			if fst.Ident == pfst.Ident {
// 				for _, selec := range selecs {
// 					if fst.MatchSelec(selec) {
// 						combos = append(combos, Combo{fst, pfst, selec})
// 						break
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return
// }

func Combine(fsts []Fst, pfsts []Fst, selecs []Fst) (combos []Combo) {
	for _, fst := range fsts {
		pfst := Fst{}
		foundpfst := false
		for _, p := range pfsts {
			if fst.Ident == p.Ident {
				pfst = p
				foundpfst = true
				break
			}
		}
		if !foundpfst {
			break
		}

		selec := Fst{}
		for _, s := range selecs {
			if fst.MatchSelec(selec) {
				selec = s
				break
			}
		}
		combos = append(combos, Combo{fst, pfst, selec})
	}
	return
}

func FprintComboFstPfstSelec(w io.Writer, combos []Combo) {
	for _, c := range combos {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", c.Pfst.Path, c.Fst.Path, c.Selec.Path, c.Fst.Ident + "_multiplot")
	}
}

func CombineFstPfstSelec() {
	all_fsts, err := GetPathFsts(os.Args[1])
	if err != nil {panic(err)}
	all_pfsts, err := GetPathPfsts(os.Args[2])
	if err != nil {panic(err)}
	all_selecs, err := GetPathSelecs(os.Args[3])
	if err != nil {panic(err)}

	combo := Combine(all_fsts, all_pfsts, all_selecs)
	FprintComboFstPfstSelec(os.Stdout, combo)
}

package main

import (
	"github.com/jgbaldwinbrown/gggenes_plot/fthresh"
)

type Fst struct {
	Ident string
	Path string
}

type Combo struct {
	Fst Fst
	Pfst Fst
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

func CombineFsts(fsts []Fst, pfsts []Fst) (combos []Combo) {
	for _, fst := range fsts {
		for _, pfst := range pfsts {
			if fst.Ident == pfst.Ident {
				combos = append(combos, Combo{fst, pfst})
			}
		}
	}
	return
}

func FprintCombo(w io.Writer, combos []Combo) {
	for _, c := range combos {
		fmt.Fprintf(w, "%v\t%v\t%v\n", c.Pfst.Path, c.Fst.Path, c.Fst.Ident + "_multiplot")
	}
}

func main() {
	all_fsts, err := GetPathFsts(os.Args[1])
	if err != nil {panic(err)}
	all_pfsts, err := GetPathPfsts(os.Args[2])
	if err != nil {panic(err)}

	combo := CombineFsts(all_fsts, all_pfsts)
	FprintCombo(os.Stdout, combo)
}



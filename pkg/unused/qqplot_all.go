package main

import (
	"os"
	"io"
	"bufio"
	"github.com/jgbaldwinbrown/makem"
	"strings"
)

type PlotSet struct {
	Pfst string
	Fst string
	Selec string
	Out string
}

type PlotSets []PlotSet

func ParsePlotSet(s string) PlotSet {
	line := strings.Split(s, "\t")
	return PlotSet{
		Pfst: line[0],
		Fst: line[1],
		Selec: line[2],
		Out: line[3],
	}
}

func ReadPlotSets(r io.Reader) (ps PlotSets) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		ps = append(ps, ParsePlotSet(s.Text()))
	}
	return
}

func FstOuts(prefix string) []string {
	return []string{prefix + "_fst_qq.pdf"}
}

func SelOuts(prefix string) []string {
	return []string{prefix + "_sel_qq.pdf"}
}

func PfstOuts(prefix string) []string {
	return []string {
		prefix + "_pfst_qq.pdf",
		prefix + "_pfst_qq_ylog.pdf",
		prefix + "_pfst_qq_nolog.pdf",
		prefix + "_pfst_qq_exp.pdf",
	}
}

func PfstNowinOuts(prefix string) []string {
	return []string {
		prefix + "_pfst_nowin_qq.pdf",
		prefix + "_pfst_nowin_qq_ylog.pdf",
		prefix + "_pfst_nowin_qq_nolog.pdf",
		prefix + "_pfst_nowin_qq_exp.pdf",
	}
}

func AddPlotSet(m *makem.MakeData, p PlotSet) {
	pfst_in := p.Out + "_pfst_plfmt.bed"
	fst_in := p.Out + "_fst_plfmt.bed"
	selec_in := p.Out + "_selec_plfmt.bed"

	fst_outs := FstOuts(p.Out)
	sel_outs := SelOuts(p.Out)
	pfst_outs := PfstOuts(p.Out)
	// pfst_nowin_outs := PfstNowinOuts(p.Out)

	r := makem.Recipe{}
	r.AddTargets(fst_outs...)
	r.AddDeps(fst_in)
	r.AddScripts("qqfst $^ " + strings.Join(fst_outs, " "))
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(sel_outs...)
	r.AddDeps(selec_in)
	r.AddScripts("qqsel $^ " + strings.Join(sel_outs, " "))
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(pfst_outs...)
	r.AddDeps(pfst_in)
	r.AddScripts("qqpfst $^ " + strings.Join(pfst_outs, " "))
	m.Add(r)

	// r = makem.Recipe{}
	// r.AddTargets(pfst_nowin_outs...)
	// r.AddDeps(pfst_nowin_in)
	// r.AddScripts("qqpfst_nowin $^ " + strings.Join(pfst_nowin_outs, " "))
	// m.Add(r)
}

func MakeMakefile(r io.Reader) *makem.MakeData {
	makefile := new(makem.MakeData)

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)

	for s.Scan() {
		AddPlotSet(makefile, ParsePlotSet(s.Text()))
	}

	return makefile
}

func main() {
	makefile := MakeMakefile(os.Stdin)
	makefile.Fprint(os.Stdout)
	err := makefile.Exec(makem.UseCores(8), makem.KeepGoing())
	if err != nil {
		panic(err)
	}
}

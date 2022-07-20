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

func AddPlotSet(m *makem.MakeData, p PlotSet) {
	pfst_plfmt := p.Out + "_pfst_plfmt.bed"
	fst_plfmt := p.Out + "_fst_plfmt.bed"
	selec_plfmt := p.Out + "_selec_plfmt.bed"
	out := p.Out + "_plot_pfst_fst_selec_v2.png"

	r := makem.Recipe{}
	r.AddTargets(pfst_plfmt)
	r.AddDeps(p.Pfst)
	r.AddScripts("./plfmt_flex -c 0 -b 2 <$< > $@")
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(fst_plfmt)
	r.AddDeps(p.Fst)
	r.AddScripts("./plfmt_flex -c 0 -b 2 <$< > $@")
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(selec_plfmt)
	r.AddDeps(p.Selec)
	r.AddScripts("./plfmt_flex -c 0 -b 1 -H <$< > $@")
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(out)
	r.AddDeps(pfst_plfmt, fst_plfmt, selec_plfmt)
	r.AddScripts("Rscript plot_pfst_fst_selec.R $^ " + out)
	m.Add(r)
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

package fthresh

import (
	"os"
	"io"
	"bufio"
	"github.com/jgbaldwinbrown/makem"
	"fmt"
)

func AddManhattanPlotSetNowin(m *makem.MakeData, p PlotSet, chrlenspath string) {
	pfst_plfmt := p.Out + "_pfst_plfmt.bed"
	fst_plfmt := p.Out + "_fst_plfmt.bed"
	selec_plfmt := p.Out + "_selec_plfmt.bed"
	out := p.Out + "_plot_pfst_fst_selec.png"
	out_noselec := p.Out + "_plot_pfst_fst.png"

	r := makem.Recipe{}
	r.AddTargets(pfst_plfmt)
	r.AddDeps(p.Pfst)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -C " + chrlenspath + " -c 0 -b 1 <$< > $@")
	} else {
		panic(fmt.Errorf("no chrlenspath"))
		r.AddScripts("plfmt_flex -c 0 -b 1 <$< > $@")
	}
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(fst_plfmt)
	r.AddDeps(p.Fst)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -C " + chrlenspath + " -c 0 -b 1 <$< > $@")
	} else {
		panic(fmt.Errorf("no chrlenspath"))
		r.AddScripts("plfmt_flex -c 0 -b 1 <$< > $@")
	}
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(selec_plfmt)
	r.AddDeps(p.Selec)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -C " + chrlenspath + " -c 0 -b 1 -H <$< > $@")
	} else {
		panic(fmt.Errorf("no chrlenspath"))
		r.AddScripts("plfmt_flex -c 0 -b 1 -H <$< > $@")
	}
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(out)
	r.AddDeps(pfst_plfmt, fst_plfmt, selec_plfmt)
	r.AddScripts("Rscript plot_pfst_fst_selec.R $^ " + out)

	r = makem.Recipe{}
	r.AddTargets(out_noselec)
	r.AddDeps(pfst_plfmt, fst_plfmt)
	r.AddScripts("Rscript plot_pfst_fst.R $^ " + out_noselec)
	m.Add(r)
}

func MakeManhatMakefileNowin(r io.Reader, chrlenpath string) *makem.MakeData {
	makefile := new(makem.MakeData)

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)

	for s.Scan() {
		AddManhattanPlotSetNowin(makefile, ParsePlotSet(s.Text()), chrlenpath)
	}

	return makefile
}

func MakeAndRunManhatMakefileNowin() {
	chrlenpath := ""
	if len(os.Args) > 1 {
		chrlenpath = os.Args[1]
	}
	if chrlenpath == "" {
		panic(fmt.Errorf("no chrlenpath at input"))
	}
	makefile := MakeManhatMakefileNowin(os.Stdin, chrlenpath)

	makefile.Fprint(os.Stdout)

	mf, err  := os.Create("manhat_makefile_nowin")
	if err != nil {
		panic(err)
	}
	makefile.Fprint(mf)
	mf.Close()

	err = makefile.Exec(makem.UseCores(8), makem.KeepGoing())
	if err != nil {
		panic(err)
	}
}


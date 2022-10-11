package fthresh

import (
	"os"
	"io"
	"bufio"
	"github.com/jgbaldwinbrown/makem"
	"strings"
	"fmt"
)

func AddManhattanPlotSet(m *makem.MakeData, p PlotSet, chrlenspath string) {
	pfst_plfmt_noslop := p.Out + "_pfst_plfmt_noslop.bed"
	pfst_plfmt := p.Out + "_pfst_plfmt.bed"
	fst_plfmt := p.Out + "_fst_plfmt.bed"
	selec_plfmt := p.Out + "_selec_plfmt.bed"
	out := p.Out + "_plot_pfst_fst_selec.png"
	out_noselec := p.Out + "_plot_pfst_fst.png"

	r := makem.Recipe{}
	r.AddTargets(pfst_plfmt_noslop)
	r.AddDeps(p.Pfst)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -C " + chrlenspath + " -c 0 -b 2 <$< > $@")
	} else {
		panic(fmt.Errorf("no chrlenspath"))
		r.AddScripts("plfmt_flex -c 0 -b 2 <$< > $@")
	}
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(pfst_plfmt)
	r.AddDeps(pfst_plfmt_noslop)
	r.AddScripts(`mawk -F "\t" -v OFS="\t" '{$$2-=24999; $$3+=25000; print $$0}' < $< > $@`)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(fst_plfmt)
	r.AddDeps(p.Fst)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -C " + chrlenspath + " -c 0 -b 2 <$< > $@")
	} else {
		panic(fmt.Errorf("no chrlenspath"))
		r.AddScripts("plfmt_flex -c 0 -b 2 <$< > $@")
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

func MakeManhatMakefileOldCfg(r io.Reader, chrlenpath string) *makem.MakeData {
	makefile := new(makem.MakeData)

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)

	for s.Scan() {
		AddManhattanPlotSet(makefile, ParsePlotSet(s.Text()), chrlenpath)
	}

	return makefile
}

func MakeManhatMakefile(r io.Reader, chrlenpath string) *makem.MakeData {
	makefile := new(makem.MakeData)

	config, err := ReadComboConfig(r)
	if err != nil {
		panic(err)
	}

	for _, cfg := range config {
		plotset := ConfigToPlotSet(cfg)
		AddManhattanPlotSet(makefile, plotset, chrlenpath)
	}

	return makefile
}

func MakeAndRunManhatMakefile() {
	chrlenpath := ""
	if len(os.Args) > 1 {
		chrlenpath = os.Args[1]
	}
	if chrlenpath == "" {
		panic(fmt.Errorf("no chrlenpath at input"))
	}
	makefile := MakeManhatMakefile(os.Stdin, chrlenpath)

	makefile.Fprint(os.Stdout)

	mf, err  := os.Create("manhat_makefile")
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

func AddQqPlotSet(m *makem.MakeData, p PlotSet) {
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

func MakeQqMakefile(r io.Reader) *makem.MakeData {
	makefile := new(makem.MakeData)

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)

	for s.Scan() {
		AddQqPlotSet(makefile, ParsePlotSet(s.Text()))
	}

	return makefile
}

func MakeAndRunQqMakefile() {
	makefile := MakeQqMakefile(os.Stdin)
	makefile.Fprint(os.Stdout)
	err := makefile.Exec(makem.UseCores(8), makem.KeepGoing())
	if err != nil {
		panic(err)
	}
}

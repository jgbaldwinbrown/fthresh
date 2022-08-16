package fthresh

import (
	"io"
	"os"
	"bufio"
	"github.com/jgbaldwinbrown/makem"
)

func AddGoodPlotSet(m *makem.MakeData, p PlotSet) {
	pfst_plfmt_noslop := p.Out + "_pfst_plfmt_noslop.bed"
	pfst_plfmt := p.Out + "_pfst_plfmt.bed"
	fst_plfmt := p.Out + "_fst_plfmt.bed"
	selec_plfmt := p.Out + "_selec_plfmt.bed"

	goodpre := p.GoodPfstSpans
	goods := goodpre
	goods_plfmt := goodpre + "_plfmt.bed"
	goods_plot := goodpre + "_plot.png"

	out := p.Out + "_plot_pfst_fst_selec.png"
	out_noselec := p.Out + "_plot_pfst_fst.png"

	r := makem.Recipe{}
	r.AddTargets(pfst_plfmt_noslop)
	r.AddDeps(p.Pfst)
	r.AddScripts("plfmt_flex -c 0 -b 2 <$< > $@")
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(pfst_plfmt)
	r.AddDeps(pfst_plfmt_noslop)
	r.AddScripts(`awk -F "\t" -v OFS="\t" '{$$2-=24999; $$3+=25000; print $$0}' < $< > $@`)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(fst_plfmt)
	r.AddDeps(p.Fst)
	r.AddScripts("plfmt_flex -c 0 -b 2 <$< > $@")
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(selec_plfmt)
	r.AddDeps(p.Selec)
	r.AddScripts("plfmt_flex -c 0 -b 1 -H <$< > $@")
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

	r = makem.Recipe{}
	r.AddTargets(goods_plfmt)
	r.AddDeps(goods)
	r.AddScripts("plfmt_flex -c 0 -b 1 -b2 2 < $< > $@")
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(goods_plot)
	r.AddDeps(pfst_plfmt, fst_plfmt, selec_plfmt, goods_plfmt)
	r.AddScripts("Rscript plot_goodss.R $^ " + goods_plot)
	m.Add(r)

	// pfst_path = args[1]
	// fst_path = args[2]
	// selec_path = args[3]
	// out_path = args[4]

	// pfst_sig_path = args[5]
	// fst_sig_path = args[6]
	// selec_sig_path = args[7]
}

func MakeGoodsMakefile(r io.Reader) *makem.MakeData {
	makefile := new(makem.MakeData)

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)

	for s.Scan() {
		AddManhattanPlotSet(makefile, ParsePlotSet(s.Text()))
	}

	return makefile
}

func RunMakefile(mf *makem.MakeData, cores int) {
	mfp, err  := os.Create("manhat_makefile")
	if err != nil {
		panic(err)
	}
	mf.Fprint(mfp)
	mfp.Close()

	err = mf.Exec(makem.UseCores(cores), makem.KeepGoing())
	if err != nil {
		panic(err)
	}
}

func MakeAndRunGoodsMakefile() {
	mf := MakeGoodsMakefile(os.Stdin)
	mf.Fprint(os.Stdout)
	RunMakefile(mf, 8)
}

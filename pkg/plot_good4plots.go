package fthresh

import (
	"io"
	"os"
	"bufio"
	"github.com/jgbaldwinbrown/makem"
)

func AddGoodPlotSet(m *makem.MakeData, p PlotSet, chrlenspath string) {
	pfst_plfmt_noslop := p.Out + "_pfst_plfmt_noslop.bed"
	pfst_plfmt := p.Out + "_pfst_plfmt.bed"
	fst_plfmt := p.Out + "_fst_plfmt.bed"
	selec_plfmt := p.Out + "_selec_plfmt.bed"

	goodpre := p.GoodPfstSpans
	goods := goodpre
	goods_plfmt := goodpre + "_plfmt.bed"
	goods_plot := goodpre + "_plot.png"

	goods_subfull := SubFullPath(goodpre)
	goods_subfull_plfmt := goods_subfull + "_plfmt.bed"
	goods_subfull_plot := goods_subfull + "_plot.bed"

	out := p.Out + "_plot_pfst_fst_selec.png"
	out_noselec := p.Out + "_plot_pfst_fst.png"

	r := makem.Recipe{}
	r.AddTargets(pfst_plfmt_noslop)
	r.AddDeps(p.Pfst)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -C " + chrlenspath + " -c 0 -b 2 <$< > $@")
	} else {
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
		r.AddScripts("plfmt_flex -c 0 -b 2 <$< > $@")
	}
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(selec_plfmt)
	r.AddDeps(p.Selec)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -C " + chrlenspath + " -c 0 -b 1 -H <$< > $@")
	} else {
		r.AddScripts("plfmt_flex -c 0 -b 1 -H <$< > $@")
	}
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(out)
	r.AddDeps(pfst_plfmt, fst_plfmt, selec_plfmt)
	r.AddScripts("#plot_pfst_fst_selec $^ " + out)

	r = makem.Recipe{}
	r.AddTargets(out_noselec)
	r.AddDeps(pfst_plfmt, fst_plfmt)
	r.AddScripts("#plot_pfst_fst $^ " + out_noselec)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(goods_plfmt)
	r.AddDeps(goods)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -c 0 -b 1 -b2 2 -C " + chrlenspath + " < $< > $@")
	} else {
		r.AddScripts("plfmt_flex -c 0 -b 1 -b2 2 < $< > $@")
	}
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(goods_plot)
	r.AddDeps(pfst_plfmt, fst_plfmt, selec_plfmt, goods_plfmt)
	r.AddScripts("#plot_goods $^ " + goods_plot)
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(goods_subfull_plfmt)
	r.AddDeps(goods_subfull)
	if chrlenspath != "" {
		r.AddScripts("plfmt_flex -c 0 -b 1 -b2 2 -C " + chrlenspath + " < $< > $@")
	} else {
		r.AddScripts("plfmt_flex -c 0 -b 1 -b2 2 < $< > $@")
	}
	m.Add(r)

	r = makem.Recipe{}
	r.AddTargets(goods_subfull_plot)
	r.AddDeps(pfst_plfmt, fst_plfmt, selec_plfmt, goods_subfull_plfmt)
	r.AddScripts("#plot_goods $^ " + goods_subfull_plot)
	m.Add(r)

	// pfst_path = args[1]
	// fst_path = args[2]
	// selec_path = args[3]
	// out_path = args[4]

	// pfst_sig_path = args[5]
	// fst_sig_path = args[6]
	// selec_sig_path = args[7]
}

func MakeGoodsMakefileOldCfg(r io.Reader, chrlenspath string) *makem.MakeData {
	makefile := new(makem.MakeData)

	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)

	for s.Scan() {
		AddGoodPlotSet(makefile, ParsePlotSet(s.Text()), chrlenspath)
	}

	return makefile
}

func MakeGoodsMakefile(r io.Reader, chrlenspath string) *makem.MakeData {
	makefile := new(makem.MakeData)
	cfgs, err := ReadComboConfig(r)
	if err != nil {
		panic(err)
	}

	for _, cfg := range cfgs {
		AddGoodPlotSet(makefile, ConfigToPlotSet(cfg), chrlenspath)
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
	mf := MakeGoodsMakefile(os.Stdin, os.Args[1])
	mf.Fprint(os.Stdout)
	RunMakefile(mf, 8)
}

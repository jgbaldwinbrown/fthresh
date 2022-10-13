package main

import (
	"github.com/jgbaldwinbrown/fthresh/pkg"
	"github.com/jgbaldwinbrown/makem"
	"os"
)

func main() {
	mf := fthresh.MakeGoodsMakefile(os.Stdin, os.Args[1])
	mf.Fprint(os.Stdout)
	err := mf.Exec(makem.UseCores(8), makem.KeepGoing())
	if err != nil {
		panic(err)
	}
}

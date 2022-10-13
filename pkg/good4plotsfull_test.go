package fthresh

import (
	"os/exec"
	"os"
	"testing"
	"io/ioutil"
)

func TestPercThreshAndMergeFull(t *testing.T) {
	inpath, err := WritePTMInput()
	if err != nil {
		panic(err)
	}
	defer os.Remove(inpath)

	outconn, err := ioutil.TempFile(".", "TestPTMOut.bed")
	if err != nil {
		panic(err)
	}
	outpath := outconn.Name()
	outconn.Close()
	defer os.Remove(outpath)
	defer os.Remove(outpath + "_thresholded.bed")
	defer os.Remove(outpath + "_thresh_merge.bed")

	err = PercThreshAndMerge(inpath, 3, 0.21, 0, outpath)
	if err != nil {
		panic(err)
	}
	cat := exec.Command("cat", outpath + "_thresh_merge.bed")
	cat.Stdout = os.Stdout
	err = cat.Run()
	if err != nil {
		panic(err)
	}
}

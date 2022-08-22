package fthresh

import (
	"os/exec"
	"fmt"
	"os"
	"testing"
	"io/ioutil"
)

var PTMInput string = `one	0	1	1
one	1	2	5
one	2	3	3
one	3	4	4
one	4	5	2
`

func WritePTMInput() (path string, err error) {
	file, err := ioutil.TempFile(".", "TestPTMIn.bed")
	if err != nil {
		return "", err
	}
	defer file.Close()
	fmt.Fprint(file, PTMInput)
	return file.Name(), nil
}

func TestPercThreshAndMerge(t *testing.T) {
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

	err = PercThreshAndMerge(inpath, 3, 0.21, outpath)
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

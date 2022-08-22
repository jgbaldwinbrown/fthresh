package main

import (
	"testing"
	"strings"
	"os"
)



func TestPlfmt(t *testing.T) {
	data := `chr1	25	50
chr1	33	98
chr2	100	101
chr3	105	107`
	r := strings.NewReader(data)
	f := Flags {
		Header: false,
		ChrCol: 0,
		BpCol: 1,
		BpCol2: 2,
		ChrBedPath: "chroms.bed",
	}
	out, err := os.Create("testout.bed")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	Plfmt(f, r, out)
}

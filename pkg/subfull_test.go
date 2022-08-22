package fthresh

import (
	"fmt"
	"testing"
)

// func SubtractAlts(gset GoodAndAlts, statistic string) (outpath string, err error) {

// func WriteSubtractTestFiles() (todelete []string, err error) {
// 	
// }

func TestSubtractAlts(t *testing.T) {
	gset := GoodAndAlts{
		Comp{
			Breed1: "black",
			Bit1: "unbitted",
			Breed2: "feral",
			Bit2: "unbitted",
		},
		[]Comp{
			Comp{
				Breed1: "black",
				Bit1: "bitted",
				Breed2: "feral",
				Bit2: "bitted",
			},
		}
	}
	
}

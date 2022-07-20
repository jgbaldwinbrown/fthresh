package main

import (
	"github.com/chewxy/stl/loess"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"fmt"
	"sort"
	"strconv"
	"errors"
	"github.com/jgbaldwinbrown/lscan/lscan"
	"bufio"
	"io"
	"os"
)

func Deriv(x []float64, y []float64) (xout, yout []float64, err error) {
	if len(x) < 1 {
		return xout, yout, errors.New("Deriv error: too short.")
	}
	if len(x) != len(y) {
		return xout, yout, errors.New("Deriv error: lengths do not match.")
	}
	xout = make([]float64, len(x) - 1)
	yout = make([]float64, len(y) - 1)
	l := len(x) - 1
	for i:=0; i<l; i++ {
		xout[i] = (x[i] + x[i+1]) / 2
		yout[i] = (y[i+1] - y[i]) / (x[i+1] - x[i])
	}
	return xout, yout, nil
}

func ReadTable(r io.Reader, col int) ([]float64, error) {
	var data []float64
	var line []string
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	splitter := lscan.ByByte('\t')
	for s.Scan() {
		line = lscan.SplitByFunc(line, s.Text(), splitter)
		f, err := strconv.ParseFloat(line[col], 64)
		if err != nil {
			return data, errors.New("Line too short.")
		}
		data = append(data, f)
	}
	sort.Float64s(data)
	return data, nil
}

func DistribIndices(n int, f func(float64) float64) []float64 {
	out := make([]float64, n)
	fn := float64(n)

	for i:=0; i<n; i++ {
		out[i] = f(float64(i)/fn)
	}
	return out
}

func ExpDistf() func(p float64) float64 {
	dist := distuv.Exponential{}
	dist.Rate = 12.5
	return func(p float64) float64 {
		return dist.Quantile(p)
	}
}

func CalcAccelFull(param string) {
	var distf func(float64)float64
	col := 3
	switch param {
	case "fst":
		distf = ExpDistf()
	case "selec":
		distf = ExpDistf()
	case "pfst":
		col = 8
		distf = ExpDistf()
	default:
		distf = func(p float64) float64 { return p }
	}

	data, err := ReadTable(os.Stdin, col)
	if err != nil {
		fmt.Println(err)
		return
	}

	y, err := loess.Smooth(data, 3000, 1, loess.Linear)
	if err != nil {
		panic(err)
	}

	x := DistribIndices(len(data), distf)

	lx := len(x)
	ly := len(y)
	x = x[lx/4:lx-(lx/20000)]
	y = y[ly/4:ly-(ly/20000)]

	xd1, yd1, _ := Deriv(y, x)
	_, yd2, _ := Deriv(xd1, yd1)

	if len(yd2) < 1 {
		fmt.Println("Empty file.")
		return
	}
	max := yd2[0]
	maxi := 0
	for i, y := range yd2[1:] {
		if math.IsNaN(max) || math.IsInf(max, 0) || (max < y && !math.IsInf(y, 0)) {
			max = y
			maxi = i
		}
	}

	fmt.Println(data[maxi], max)
}

func main() {
	param := os.Args[1]
	CalcAccelFull(param)
}

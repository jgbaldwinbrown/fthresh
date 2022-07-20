package main

import (
	"os"
	"flag"
	"fmt"
	"strconv"
	"io"
	"bufio"
	"strings"
	"sort"
	//"github.com/pkg/profile"
)

type Flags struct {
	Header bool
	ChrCol int
	BpCol int
}

type Entry struct {
	Line []string
	Chr int
	Bp int
}

type Entries []Entry
func (e Entries) Len() int { return len(e) }
func (e Entries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Entries) Less(i, j int) bool {
	if e[i].Chr == e[j].Chr {
		return e[i].Bp < e[j].Bp
	}
	return e[i].Chr < e[j].Chr
}

func GetEntry(line []string, chrcol int, bpcol int) (e Entry, err error) {
	e.Line = line

	chrnum, err := strconv.ParseFloat(line[chrcol][3:], 64)
	if err != nil { return }
	e.Chr = int(chrnum)

	bp, err := strconv.ParseFloat(line[bpcol], 64)
	e.Bp = int(bp)
	return
}

func GetData(r io.Reader, chrcol int, bpcol int, header bool) (es Entries, err error, header_string string) {
	s := bufio.NewScanner(r)

	if header {
		s.Scan()
		header_string = s.Text()
	}

	for s.Scan() {
		line := strings.Split(s.Text(), "\t")
		var e Entry
		e, err = GetEntry(line, chrcol, bpcol)
		if err != nil { return }
		es = append(es, e)
	}
	return
}

// func Flush(buf *Buf) {
// 	lines := strings.Join(buf.Lines, "\n")
// 	fmt.Fprintf(buf.Writer, "%s\n", lines)
// 	buf.Lines = buf.Lines[:0]
// }

func FprintEntry(w io.Writer, e Entry, cumsum int) {
	fmt.Fprintf(w, "%s\t%v\t%v\n", strings.Join(e.Line, "\t"), e.Chr, cumsum)
	// line := fmt.Sprintf("%s\t%v\t%v\n", strings.Join(e.Line, "\t"), e.Chr, cumsum)
	// buf.Lines = append(buf.Lines, line)
	// if len(buf.Lines) > buf.Max {
	// 	Flush(buf)
	// }
}

// type Buf struct {
// 	Writer io.Writer
// 	Lines []string
// 	Max int
// }

func append_pls(data Entries, w io.Writer) {
	sort.Sort(data)
	cumsum := 0
	prevchr := -1
	prevbp := -1

	// outbuf := Buf{}
	// outbuf.Max = 10000
	// outbuf.Writer = w
	// defer Flush(&outbuf)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for _, e := range data {
		if e.Chr != prevchr {
			cumsum += 1000
		} else {
			cumsum += e.Bp - prevbp
		}
		FprintEntry(bw, e, cumsum)
		prevchr = e.Chr
		prevbp = e.Bp
	}
}

func Plfmt(flags Flags, r io.Reader, w io.Writer) {
	data, err, _ := GetData(r, flags.ChrCol, flags.BpCol, flags.Header)
	if err != nil { panic(err) }
	append_pls(data, w)
}

func GetFlags() (f Flags) {
	flag.BoolVar(&f.Header, "H", false, "File includes a header line.")
	flag.IntVar(&f.ChrCol, "c", -1, "0-indexed column containing chromosome in format \"chr[0-9]*\"")
	flag.IntVar(&f.BpCol, "b", -1, "Column containing basepair position in format \"chr[0-9]*\"")
	flag.Parse()
	if f.ChrCol == -1 || f.BpCol == -1 {
		panic(fmt.Errorf("Missing argument"))
	}
	return
}

func main() {
	//defer profile.Start().Stop()
	flags := GetFlags()
	Plfmt(flags, os.Stdin, os.Stdout)
}

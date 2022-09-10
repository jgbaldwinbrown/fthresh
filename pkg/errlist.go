package fthresh

import (
	"strings"
	"fmt"
)

type errList []error

func (l errList) Error() string {
	var b strings.Builder
	fmt.Fprint(&b, "[")
	if len(l) > 0 {
		fmt.Fprint(&b, l[0].Error())
	}
	for _, e := range l[1:] {
		fmt.Fprintf(&b, ", %s", e.Error())
	}
	fmt.Fprint(&b, "]")
	return b.String()
}

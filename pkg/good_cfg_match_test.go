package fthresh

import (
	"fmt"
	"testing"
)

func CompEx() Comp {
	return Comp{
		"black",
		"bitted",
		2,
		36,
		"feral",
		"bitted",
		2,
		36,
	}
}

func CompExAll() Comp {
	return Comp{
		"black",
		"bitted",
		0,
		36,
		"feral",
		"bitted",
		0,
		36,
	}
}

func ComboEx() ComboConfig {
	return ComboConfig{
		Treatment: Treatment{
			"black",
			"bitted",
			"36",
			"R2",
		},
		Treatment2: Treatment{
			"feral",
			"bitted",
			"36",
			"R2",
		},
	}
}

func ComboExAll() ComboConfig {
	return ComboConfig{
		Treatment: Treatment{
			"black",
			"bitted",
			"36",
			"All",
		},
		Treatment2: Treatment{
			"feral",
			"bitted",
			"36",
			"All",
		},
	}
}


func TestMatchConfig(t *testing.T) {
	if !MatchConfigAndGoodAndAlts(ComboEx(), CompEx()) {
		t.Errorf("combo %v and comp %v don't match", ComboEx(), CompEx())
	}
}

func TestMatchConfigAll(t *testing.T) {
	if !MatchConfigAndGoodAndAlts(ComboExAll(), CompExAll()) {
		t.Errorf("combo %v and comp %v don't match", ComboExAll(), CompExAll())
	}
}

func TestFindMatch(t *testing.T) {
	idx := FindMatchingConfig(CompEx(), []ComboConfig{ComboEx()})
	if idx == -1 {
		t.Errorf("no match")
	}
	fmt.Println("match")
}

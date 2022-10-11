package fthresh

import (
	"os"
	"io"
	"encoding/json"
)

type Treatment struct {
	Breed string `json:"breed"`
	Bit string `json:"bit"`
	Time string `json:"time"`
	Replicate string `json:"replicate"`
}

type ComboConfig struct {
	Treatment Treatment `json:"treatment"`
	Treatment2 Treatment `json:"treatment2"`
	Pfst string `json:"pfst"`
	Fst string `json:"fst"`
	Selec string `json:"selec"`
	WinSize string `json:"winsize"`
	WinStep string `json:"winstep"`
	ComparisonType string `json:"comparison_type"`
	OutPrefix string `json:"out_prefix"`
	Subtractions string `json:"subtractions"`
}

func ReadComboConfig(r io.Reader) ([]ComboConfig, error) {
	cfgbytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var cfg []ComboConfig
	err = json.Unmarshal(cfgbytes, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetComboConfig(path string) ([]ComboConfig, error) {
	if path == "" {
		return ReadComboConfig(os.Stdin)
	}

	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	return ReadComboConfig(r)
}

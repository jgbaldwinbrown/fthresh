package fthresh

import (
	"strings"
	"fmt"
)

func SubFullPath(path string) string {
	return strings.Replace(path, ".bed", "_subfulls.bed", 1)
}

func BreedPriority() map[string]int {
	return map[string]int {
		"black": 5,
		"runt": 4,
		"white": 3,
		"figurita": 2,
		"feral": 1,
	}
}

func BitPriority() map[string]int {
	return  map[string]int {
		"unbitted": 2,
		"bitted": 1,
	}
}

func FormatMap() map[string]string {
	return map[string]string {
		"black": "BlackHomer",
		"white": "WhiteHomer",
		"feral": "Feral",
		"runt": "Runt",
		"figurita": "Figurita",
		"bitted": "Bitted",
		"unbitted": "Unbitted",
	}
}

func HeightMap() map[string]string {
	return  map[string]string {
		"black": "Low",
		"white": "Low",
		"feral": "Mid",
		"figurita": "Low",
		"runt": "Low",
	}
}

// func HeightMap() map[string]string {
// 	return  map[string]string {
// 		"black": "Low",
// 		"white": "High",
// 		"feral": "Mid",
// 		"figurita": "Low",
// 		"runt": "High",
// 	}
// }

func TraitMap() map[string]string {
	return  map[string]string {
		"black": "Color",
		"white": "Color",
		"feral": "Color",
		"figurita": "Size",
		"runt": "Size",
	}
}

func AltMap() map[string][][]string {
	return  map[string][][]string {
		"black": [][]string{
			[]string{ "white", "unbitted" },
			[]string{ "feral", "unbitted" },
			[]string{ "black", "bitted" },
		},
		"white": [][]string{
			[]string{ "white", "bitted" },
			[]string{ "feral", "unbitted" },
			[]string{ "black", "unbitted" },
		},
		"figurita": [][]string{
			[]string{ "runt", "unbitted" },
			[]string{ "feral", "unbitted" },
			[]string{ "figurita", "bitted" },
		},
		"runt": [][]string{
			[]string{ "runt", "bitted" },
			[]string{ "feral", "unbitted" },
			[]string{ "figurita", "unbitted" },
		},
		"feral": [][]string{
			[]string{ "feral", "bitted" },
			[]string{ "black", "unbitted" },
			[]string{ "white", "unbitted" },
			[]string{ "figurita", "unbitted" },
			[]string{ "runt", "unbitted" },
		},
	}
}

// func BitMap() map[string]string {
// 	return map[string]string {
// 		"black": "BlackHomer",
// 		"white": "WhiteHomer",
// 		"feral": "Feral",
// 		"figurita": "Figurita",
// 		"runt": "Runt",
// 	}
// }

func reformatStr(str string) string {
	return ConservativeMap(str, FormatMap())
}

func ConservativeMap[T comparable](str T, m map[T]T) T {
	out, ok := m[str]
	if !ok {
		return str
	}
	return out
}

func SyncString(breed, bit string) string {
	return fmt.Sprintf("split_all_pools/sync/sync2/%v_tall_names_%v_repall.split_goods_f0.1_c10.sync.gz", breed, bit)
}

func GggenesBedString(breed1, bit1, breed2, bit2 string) string {
	return fmt.Sprintf(
		"_breed_%v_time_36_bit_%v_replicate_All_breed_%v_time_36_bit_%v_replicate_All_%v_%v_%v___multiplot_gggenes_thresholded_intervals.bed",
		FormatMap()[breed1],
		FormatMap()[bit1],
		FormatMap()[breed2],
		FormatMap()[bit2],
		TraitMap()[breed1],
		HeightMap()[breed1],
		HeightMap()[breed2],
	)
}

func PlfmtMap() map[string]string {
	return map[string]string {
		"pFst": "pfst_plfmt",
		"Fst": "fst_plfmt",
		"Selec": "selec_plfmt_bedified",
	}
}

func BedString(breed1, bit1, breed2, bit2, statistic string) string {
	return fmt.Sprintf(
		"_breed_%v_time_36_bit_%v_replicate_All_breed_%v_time_36_bit_%v_replicate_All_%v_%v_%v___multiplot_%v_tm_perc_thresh_merge.bed",
		FormatMap()[breed1],
		FormatMap()[bit1],
		FormatMap()[breed2],
		FormatMap()[bit2],
		TraitMap()[breed1],
		HeightMap()[breed1],
		HeightMap()[breed2],
		PlfmtMap()[statistic],
	)
}
// _breed_Figurita_time_36_bit_Unbitted_replicate_All_breed_Runt_time_36_bit_Unbitted_replicate_All_Size_Low_High___multiplot_fst_plfmt_tm_thresh_merge.bed
// _breed_Figurita_time_36_bit_Unbitted_replicate_All_breed_Runt_time_36_bit_Unbitted_replicate_All_Size_Low_High___multiplot_pfst_plfmt_tm_thresh_merge.bed
// _breed_Figurita_time_36_bit_Unbitted_replicate_All_breed_Runt_time_36_bit_Unbitted_replicate_All_Size_Low_High___multiplot_selec_plfmt_bedified_tm_thresh_merge.bed


func InfoString(breed, bit string) string {
	return fmt.Sprintf("split_all_pools/sync/sync2/%v_pooled_info_ne_%v.txt", breed, bit)
}

func OutprefixString(breed1, bit1, breed2, bit2 string) string {
	return fmt.Sprintf("%v_%v_%v_%v_separated", breed1, bit1, breed2, bit2)
}

func CompOutput(breed1, bit1, statistic string) string {
	return fmt.Sprintf("%v_%v_%v_subtractedalts", breed1, bit1, statistic)
}

func PrintLine(breed1, bit1, breed2, bit2 string) {
	bp1, bp2 := BreedPriority()[breed1], BreedPriority()[breed2]
	if bp2 > bp1 {
		breed1, bit1, breed2, bit2 = breed2, bit2, breed1, bit1
	}
	if bp2 == bp1 {
		bp1, bp2 = BitPriority()[bit1], BitPriority()[bit2]
		if bp2 > bp1 {
			breed1, bit1, breed2, bit2 = breed2, bit2, breed1, bit1
		}
	}
	fmt.Printf(
		"%v\t%v\t%v\t%v\t%v\n",
		SyncString(breed1, bit1),
		GggenesBedString(breed1, bit1, breed2, bit2),
		InfoString(breed1, bit1),
		OutprefixString(breed1, bit1, breed2, bit2),
		OutprefixString(breed1, bit1, breed2, bit2),
	)
}

func PrintSyncBedAndInfo() {
	breeds := []string{"black", "white", "figurita", "runt", "feral"}
	bits := []string{"bitted", "unbitted"}
	for _, breed := range breeds {
		for _, bit := range bits {
			alts := AltMap()[breed]
			for _, alt := range alts {
				PrintLine(breed, bit, alt[0], alt[1])
			}
		}
	}
}
// _breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_BlackHomer_time_36_bit_Bitted_replicate_All_Color_Low_Mid___multiplot_gggenes_thresholded_intervals.bed


// _breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Low_Mid___multiplot_fst_plfmt_tm_thresh_merge.bed
// _breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Low_Mid___multiplot_pfst_plfmt_tm_thresh_merge.bed
// _breed_WhiteHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Low_Mid___multiplot_selec_plfmt_bedified_tm_thresh_merge.bed

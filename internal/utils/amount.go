package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

func ParseAmount(amount string) float64 {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return 0
	}
	return f
}

func TruncateAmount(amount float64, n int) float64 {
	shifter := math.Pow(10, float64(n))
	return math.Floor(amount*shifter) / shifter
}

func PrintStruct(v interface{}) {
	pretty, _ := json.MarshalIndent(v, "", "\t")
	fmt.Println("\n", string(pretty))
}

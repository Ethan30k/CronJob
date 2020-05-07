package libs

import (
	"fmt"
	"math"
)

func SizeFormat(s int) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	mod := 1024.0
	i := 0
	var newnum float64 = float64(s)
	for newnum >= mod {
		newnum /= mod
		i++
	}
	return fmt.Sprintf("%.0f", math.Round(newnum)) + sizes[i]
}

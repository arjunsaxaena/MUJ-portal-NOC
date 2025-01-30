package util

import (
	"fmt"
	"strconv"
)

func FormatFloat(value string) string {
	if value == "" {
		return "0.00"
	}
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return fmt.Sprintf("%.2f", floatVal)
	}
	return "0.00"
}

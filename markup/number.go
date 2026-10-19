package markup

import (
	"math"
	"strconv"
	"strings"
)

func jsNumber(value float64) string {
	if value == 0 {
		return "0"
	}
	if math.IsNaN(value) {
		return "NaN"
	}
	if math.IsInf(value, 1) {
		return "Infinity"
	}
	if math.IsInf(value, -1) {
		return "-Infinity"
	}

	negative := value < 0
	mantissa, exponentText, _ := strings.Cut(strconv.FormatFloat(math.Abs(value), 'e', -1, 64), "e")
	digits := strings.Replace(mantissa, ".", "", 1)
	exponent, _ := strconv.Atoi(exponentText)
	point := exponent + 1
	count := len(digits)

	var out string
	switch {
	case count <= point && point <= 21:
		out = digits + strings.Repeat("0", point-count)
	case 0 < point && point <= 21:
		out = digits[:point] + "." + digits[point:]
	case -6 < point && point <= 0:
		out = "0." + strings.Repeat("0", -point) + digits
	default:
		sign := "+"
		power := point - 1
		if power < 0 {
			sign = "-"
			power = -power
		}
		if count == 1 {
			out = digits + "e" + sign + strconv.Itoa(power)
		} else {
			out = digits[:1] + "." + digits[1:] + "e" + sign + strconv.Itoa(power)
		}
	}

	if negative {
		return "-" + out
	}
	return out
}

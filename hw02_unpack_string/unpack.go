package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const EscSmb = '\\'

func Unpack(input string) (string, error) {
	// Place your code here.
	sb := strings.Builder{}
	rr := []rune(input)
	var r2Rpt rune

	for i := 0; i < len(rr); i++ {
		switch {
		case unicode.IsDigit(rr[i]):
			switch {
			case r2Rpt != 0:
				cnt, _ := strconv.Atoi(string(rr[i]))
				sb.WriteString(strings.Repeat(string(r2Rpt), cnt))
				r2Rpt = 0
			case i > 0 && rr[i-1] == EscSmb:
				r2Rpt = rr[i]
			default:
				return "", ErrInvalidString
			}
		case rr[i] == EscSmb:
			if r2Rpt != 0 {
				sb.WriteRune(r2Rpt)
				r2Rpt = 0
			} else if i > 0 && rr[i-1] == EscSmb {
				r2Rpt = rr[i]
			}
		default:
			if r2Rpt != 0 {
				sb.WriteRune(r2Rpt)
			} else if i > 0 && rr[i-1] == EscSmb {
				return "", ErrInvalidString
			}
			r2Rpt = rr[i]
		}
	}
	if r2Rpt != 0 {
		sb.WriteRune(r2Rpt)
	}
	return sb.String(), nil
}

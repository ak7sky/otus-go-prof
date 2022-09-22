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
	var r2Rpt rune
	var prevR rune

	for _, r := range input {
		switch {
		case unicode.IsDigit(r):
			switch {
			case r2Rpt != 0:
				cnt, err := strconv.Atoi(string(r))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(r2Rpt), cnt))
				r2Rpt = 0
			case prevR == EscSmb:
				r2Rpt = r
			default:
				return "", ErrInvalidString
			}
		case r == EscSmb:
			if r2Rpt != 0 {
				sb.WriteRune(r2Rpt)
				r2Rpt = 0
			} else if prevR == EscSmb {
				r2Rpt = r
			}
		default:
			if r2Rpt != 0 {
				sb.WriteRune(r2Rpt)
			} else if prevR == EscSmb {
				return "", ErrInvalidString
			}
			r2Rpt = r
		}
		prevR = r
	}

	if r2Rpt != 0 {
		sb.WriteRune(r2Rpt)
	}
	return sb.String(), nil
}

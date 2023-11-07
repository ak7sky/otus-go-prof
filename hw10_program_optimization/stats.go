package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

//go:generate easyjson stats.go
//easyjson:json
type UserEmail struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	br := bufio.NewReader(r)
	var ue UserEmail
	result := make(DomainStat)

	for {
		bytes, _, err := br.ReadLine()
		if errors.Is(err, io.EOF) {
			break
		}
		if err = ue.UnmarshalJSON(bytes); err != nil {
			return result, fmt.Errorf("err on json unmarshaling: %w", err)
		}
		if !strings.HasSuffix(ue.Email, "."+domain) {
			continue
		}
		lvl2dom := strings.ToLower(strings.SplitN(ue.Email, "@", 2)[1])
		num := result[lvl2dom]
		num++
		result[lvl2dom] = num
	}

	return result, nil
}

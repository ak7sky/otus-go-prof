package hw10programoptimization

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Old impl.
func getDomainStatOld(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

// Part of old impl.
type users [100_000]user

// Part of old impl: unmarshalling all fields, but only the "Email" field is needed.
type user struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

// Part of old impl: redundant mem alloc for all content and lines, standard json unmarshalling.
func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var u user
		if err = json.Unmarshal([]byte(line), &u); err != nil {
			return
		}
		result[i] = u
	}
	return
}

// Part of old impl: regexp using, duplicated invoke of email value splitting (minor).
func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}

const searchedDomain = "biz"

func BenchmarkGetDomainStatOld(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	defer func(r *zip.ReadCloser) { require.NoError(b, r.Close()) }(r)
	require.NoError(b, err)
	file, err := r.File[0].Open()
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = getDomainStatOld(file, searchedDomain)
	}
	b.StopTimer()
}

func BenchmarkGetDomainStatNew(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	defer func(r *zip.ReadCloser) { require.NoError(b, r.Close()) }(r)
	require.NoError(b, err)
	file, err := r.File[0].Open()
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(file, searchedDomain)
	}
	b.StopTimer()
}

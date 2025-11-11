package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetPages(page int, pageCount int) []string {
	pageStart := page - 2
	pageEnd := page + 3
	var s []string

	// append works on nil slices.

	if page < 3 {
		pageStart = 1
	}
	if page+3 > pageCount {
		pageEnd = pageCount
	}

	for i := pageStart; i <= pageEnd; i++ {
		s = append(s, fmt.Sprint(i))
	}
	return s
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func MustGetInt(key string, def int) int {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		return def
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		log.Fatalf("invalid int in %s: %q (%v)", key, v, err)
	}
	return n
}

func OnOffToBool(value string) bool {
	return value == "on"
}

func BuildPageList(total, current, maxButtons int) []int {
	if total <= 1 {
		return nil
	}
	if maxButtons < 5 {
		maxButtons = 5
	}

	list := []int{}
	if total <= maxButtons {
		for i := 1; i <= total; i++ {
			list = append(list, i)
		}
		return list
	}

	list = append(list, 1)
	left, right := current-2, current+2

	if left <= 2 {
		left, right = 2, maxButtons-1
	}
	if right >= total-1 {
		left, right = total-maxButtons+2, total-1
	}

	if left > 2 {
		list = append(list, -1)
	}

	for i := left; i <= right; i++ {
		list = append(list, i)
	}

	if right < total-1 {
		list = append(list, -1)
	}
	list = append(list, total)
	return list
}

// SplitByRune splits the string `s` by the rune separator `sep` and
// returns the resulting slices. This is a small utility used by
// path and other simple parsers that need rune-level splitting.
func SplitByRune(s string, sep rune) []string {
	var res []string
	curr := ""
	for _, c := range s {
		if c == sep {
			res = append(res, curr)
			curr = ""
		} else {
			curr += string(c)
		}
	}
	res = append(res, curr)
	return res
}

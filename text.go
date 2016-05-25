package main

import (
	"bytes"
	"strings"
	"unicode"
)

func Sanitize(r rune) rune {
	switch {
	case unicode.IsPunct(r):
		return ' '
	case unicode.IsMark(r):
		return ' '
	case unicode.IsSymbol(r):
		return ' '
	}
	return r
}

func Truncate(s string, limit int) string {
	var buf bytes.Buffer
	buf.WriteString(s)

	if buf.Len() > limit-3 {
		buf.Truncate(limit - 3)
		return buf.String() + "..."
	} else {
		return buf.String()
	}
}

func GetQueryTerms(text string) []string {
	return strings.Fields(strings.Map(Sanitize, strings.ToLower(text)))
}

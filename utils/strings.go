package utils

import "strings"

const trimCutset = "\t\n "

func TrimStrWhitespace(str string) string {
	return strings.Trim(str, trimCutset)
}

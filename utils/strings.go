package utils

import "strings"

const trimCutset = "\t\n "

func Trim(str string) string {
	return strings.Trim(str, trimCutset)
}

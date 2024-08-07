package helper

import (
	"strings"
	"unicode/utf8"
)

type MakeColumnTextWithIdentSpec struct {
	Items         []string
	Delimiter     string
	MaxRowLen     int
	IdentLen      int
	FirstRowIdent int
}

func GetTextFromItemsWithIdent(spec MakeColumnTextWithIdentSpec) string {
	builder := strings.Builder{}
	builder.WriteString(strings.Repeat(" ", spec.FirstRowIdent))

	rowLen := spec.IdentLen
	for i, item := range spec.Items {
		itemText := item
		if i < len(spec.Items)-1 {
			itemText += spec.Delimiter
		}

		if rowLen+len(itemText) > spec.MaxRowLen {
			builder.WriteString("\n")
			builder.WriteString(strings.Repeat(" ", spec.IdentLen))
			rowLen = spec.IdentLen
		}

		builder.WriteString(itemText)
		rowLen += utf8.RuneCountInString(itemText)
	}

	return builder.String()
}

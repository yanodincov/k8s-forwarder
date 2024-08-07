package forms

import (
	"github.com/manifoldco/promptui"
	"strings"
)

func GetTableInfoWithTitle(title string, kvs ...string) string {
	res := strings.Builder{}

	res.WriteString(promptui.Styler(promptui.FGBold)(title))
	res.WriteString("\n")

	// Найти максимальную длину ключа
	maxKeyLength := 0
	for i := 0; i < len(kvs); i += 2 {
		if len(kvs[i]) > maxKeyLength {
			maxKeyLength = len(kvs[i])
		}
	}

	for i := 0; i < len(kvs); i += 2 {
		key := kvs[i]
		value := kvs[i+1]

		res.WriteString(promptui.Styler(promptui.FGItalic)(key))
		res.WriteString(strings.Repeat(" ", maxKeyLength-len(key)))
		res.WriteString("  ")
		res.WriteString(value)

		if i+2 < len(kvs) {
			res.WriteString("\n")
		}
	}

	return res.String()
}

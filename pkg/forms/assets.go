package forms

import "github.com/manifoldco/promptui"

var (
	CancelName = promptui.Styler(promptui.FGItalic)("Cancel")

	CancelOptionSpec = OptionSpec{
		Data: OptionData{
			ID:          "cancel",
			Name:        CancelName,
			Description: "Return back to the previous screen",
		},
		Func: func(data OptionData) bool {
			return true
		},
	}
)

func GetCreateOptionSpec(itemName string, fn func(data OptionData) bool) OptionSpec {
	return OptionSpec{
		Data: OptionData{
			ID:          "add",
			Name:        promptui.Styler(promptui.FGItalic)("Add " + itemName),
			Description: "Add new " + itemName,
		},
		Func: fn,
	}
}

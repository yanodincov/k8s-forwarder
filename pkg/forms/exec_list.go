package forms

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/pkg/cli"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
)

type SelectFormSpecFactory func() (*SelectFormSpec, error)

type SelectFormSpec struct {
	HeaderFn   string
	QuestionFn string
	ErrorText  string
	Items      []OptionSpec
}

type OptionSpec struct {
	Data OptionData
	Func func(data OptionData) bool
}

type OptionData struct {
	ID          string
	Name        string
	Description string
}

func RunSelectForm(specFn SelectFormSpecFactory) error {
	for {
		shouldExit, err := ExecSelectForm(specFn)
		if err != nil {
			return err
		}
		if shouldExit {
			return nil
		}
	}
}

func ExecSelectForm(specFn SelectFormSpecFactory) (bool, error) {
	cli.ClearScreen()

	spec, err := specFn()
	if err != nil {
		return true, err
	}

	fnByIdx := make(map[int]func(data OptionData) bool, len(spec.Items))
	for idx, item := range spec.Items {
		fnByIdx[idx] = item.Func
	}

	items := make([]OptionData, 0, len(spec.Items))
	for _, item := range spec.Items {
		items = append(items, item.Data)
	}

	errText := " "
	if spec.ErrorText != "" {
		errText = promptui.Styler(promptui.FGRed)(spec.ErrorText)
	}

	if spec.HeaderFn != "" {
		fmt.Println("\n" + spec.HeaderFn)
	}

	selectUI := promptui.Select{
		Label: spec.QuestionFn,
		Items: items,
		Templates: &promptui.SelectTemplates{
			Help:     errText,
			Label:    "{{ .Name }}?",
			Active:   promptui.IconSelect + " {{ .Name | magenta }}",
			Inactive: promptui.Styler(promptui.FGFaint)("  {{ .Name }}"),
			Details:  "{{ .Description }}",
		},
	}

	idx, _, err := selectUI.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return true, nil
		}
		return false, err
	}

	if len(fnByIdx) > idx {
		return fnByIdx[idx](items[idx]), nil
	}

	return false, nil
}

func CreatePromptUiSelect(spec SelectFormSpec) promptui.Select {
	return promptui.Select{
		Label: spec.QuestionFn,
		Items: helper.SliceMap(spec.Items, func(item OptionSpec) OptionData {
			return item.Data
		}),
		Templates: &promptui.SelectTemplates{
			Help:     spec.ErrorText,
			Label:    "{{ .ServiceName }}?",
			Active:   promptui.IconSelect + " {{ .ServiceName | magenta }}",
			Inactive: promptui.Styler(promptui.FGFaint)("  {{ .ServiceName }}"),
			Details:  "{{ .Description }}",
		},
	}

}

package forms

import (
	"github.com/AlecAivazis/survey/v2"
	"strings"
)

func FixSurveyStyle() {
	survey.ConfirmQuestionTemplate = strings.ReplaceAll(survey.ConfirmQuestionTemplate, "cyan", "magenta")
	survey.InputQuestionTemplate = strings.ReplaceAll(survey.InputQuestionTemplate, "cyan", "magenta")
	survey.MultiSelectQuestionTemplate = strings.ReplaceAll(survey.MultiSelectQuestionTemplate, "cyan", "magenta")
	survey.PasswordQuestionTemplate = strings.ReplaceAll(survey.PasswordQuestionTemplate, "cyan", "magenta")
	survey.SelectQuestionTemplate = strings.ReplaceAll(survey.SelectQuestionTemplate, "cyan", "magenta")
	survey.EditorQuestionTemplate = strings.ReplaceAll(survey.EditorQuestionTemplate, "cyan", "magenta")
	survey.MultilineQuestionTemplate = strings.ReplaceAll(survey.MultilineQuestionTemplate, "cyan", "magenta")
	survey.ConfirmQuestionTemplate = strings.ReplaceAll(survey.ConfirmQuestionTemplate, "cyan", "magenta")
	survey.ErrorTemplate = `{{color .Icon.Format }}{{ .Icon.Text }} Sorry, your reply was invalid: {{ .Error.Error }}{{color "reset"}}
`
}

func IsInterruptSurveyErr(err error) bool {
	return strings.Contains(err.Error(), "interrupt")
}

func GetSurveyOpts() []survey.AskOpt {
	return []survey.AskOpt{
		survey.WithIcons(func(icons *survey.IconSet) {
			icons.Question.Text = ""
			icons.Question.Format = ""
			icons.Error.Text = "Error:"
			icons.SelectFocus.Text = ">"
			icons.SelectFocus.Format = "magenta"
		}),
	}
}

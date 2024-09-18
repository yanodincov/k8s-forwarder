package selectpkg

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanodincov/k8s-forwarder/pkg/forms/assets/selectpkg/template"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
)

const (
	pageSize           = 3
	defaultForeground  = "#bbbbbb"
	selectedForeground = "#ffffff"
	accentColor        = "#ff8700"
	cursorSymbol       = "→"
	pageSymbol         = " ▪"
)

func (m *Model) View() string {
	tea.ClearScreen()

	opt, i := helper.SliceFilterOne(m.currentVariants, func(opt *variant) bool {
		return opt.IsHovered
	})

	totalPagesI := len(m.currentVariants)/pageSize + 1
	currentPageI := i / pageSize

	start := currentPageI * pageSize
	end := min(start+pageSize, len(m.currentVariants))
	variants := m.currentVariants[start:end]

	view, err := template.RenderSelectTemplate(template.SelectTemplateData{
		Header:   m.headerFn(),
		Question: m.questionFn(),
		Filter:   string(m.filter),
		Footer:   opt.Option.Desc,

		TotalPages:  totalPagesI,
		CurrentPage: currentPageI,
		PageSymbol:  pageSymbol,

		HoverSymbol:   cursorSymbol,
		AccentColor:   accentColor,
		ActiveColor:   selectedForeground,
		InactiveColor: defaultForeground,

		Options: helper.SliceMap(variants, func(variant *variant) template.SelectTemplateOption {
			return template.SelectTemplateOption{
				Text:       variant.Option.Text,
				IsSelected: variant.IsSelected,
				IsHovered:  variant.IsHovered,
				IsSpecial:  !variant.QuitType.IsEmpty(),
			}
		}),
	})
	if err != nil {
		return err.Error()
	}

	return view
}

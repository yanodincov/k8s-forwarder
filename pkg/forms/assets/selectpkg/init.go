package selectpkg

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	"slices"
	"time"
)

type TickMsg struct{}

func (m *Model) Init() tea.Cmd {
	tea.ClearScreen()

	m.tickCmdFn = helper.If(m.reloadInterval > 0,
		func() tea.Cmd { return tea.Tick(m.reloadInterval, func(time.Time) tea.Msg { return TickMsg{} }) },
		func() tea.Cmd { return nil },
	)
	if m.headerFn == nil {
		m.headerFn = func() string { return "" }
	}
	if m.questionFn == nil {
		m.questionFn = func() string { return "" }
	}

	m.allVariants = append(
		helper.SliceMap(m.opts, func(opts Option) *variant {
			return &variant{Option: opts}
		}),
		&variant{Option: Option{Text: m.submitChoiceName}, QuitType: QuitTypeSubmit},
		&variant{Option: Option{Text: m.quitChoiceName}, QuitType: QuitTypeBack},
	)
	m.allVariants[0].IsHovered = true

	m.filter = []rune{}
	m.currentVariants = slices.Clone(m.allVariants)

	return m.tickCmdFn()
}

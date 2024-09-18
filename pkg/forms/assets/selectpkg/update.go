package selectpkg

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	"slices"
	"strings"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case TickMsg:
		return m, m.tickCmdFn()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			// Select the quit variant
			m.quitType = QuitTypeInterrupt

			return m, tea.Quit

		case tea.KeyUp, tea.KeyDown:
			opt, i := helper.SliceFilterOne(m.currentVariants, func(opt *variant) bool {
				return opt.IsHovered
			})
			opt.IsHovered = false

			i = helper.If(msg.Type == tea.KeyDown, i+1, i-1)
			i = min(max(i, 0), len(m.currentVariants)-1)
			m.currentVariants[i].IsHovered = true

			return m, nil

		case tea.KeyEnter:
			opt, _ := helper.SliceFilterOne(m.currentVariants, func(opt *variant) bool {
				return opt.IsHovered
			})

			// If the variant is selected, we should not allow to select special allVariants
			if !opt.QuitType.IsEmpty() {
				m.quitType = opt.QuitType

				return m, tea.Quit
			}

			opt.IsSelected = !opt.IsSelected

			return m, nil

		case tea.KeyBackspace:
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.updateVisibleOption(false)
			}

		case tea.KeyRunes, tea.KeySpace:
			m.filter = append(m.filter, msg.Runes...)
			m.updateVisibleOption(true)
		}
	}

	// No changes
	return m, nil
}

func (m *Model) updateVisibleOption(add bool) {
	if len(m.filter) == 0 {
		m.currentVariants = slices.Clone(m.allVariants)
	}

	filterStr := strings.ToLower(string(m.filter))

	var opts []*variant
	hasHovered := true
	for _, opt := range helper.If(add, m.currentVariants, m.allVariants) {

		contains := strings.Contains(strings.ToLower(opt.Option.Text), filterStr)
		if contains || !opt.QuitType.IsEmpty() {
			opts = append(opts, opt)
		}

		// If the hovered option is not in the visible options, we should remove the hover
		if opt.IsHovered && (!contains || !opt.QuitType.IsEmpty()) {
			opt.IsHovered = false
			hasHovered = false
		}
	}

	if !hasHovered && len(opts) > 0 {
		opts[0].IsHovered = true
	}

	m.currentVariants = opts
}

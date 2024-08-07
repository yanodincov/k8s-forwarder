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

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			// Select the quit innerOption
			m.quitType = QuitTypeInterrupt

			return m, tea.Quit

		case tea.KeyUp, tea.KeyDown:
			m.innerOptsQueue.Current().IsHovered = false
			helper.IfFn(msg.Type == tea.KeyDown, m.innerOptsQueue.Next, m.innerOptsQueue.Prev).IsHovered = true

			return m, nil

		case tea.KeyEnter:
			cur := m.innerOptsQueue.Current()

			// If the innerOption is selected, we should not allow to select special innerOpts
			if !cur.QuitType.IsEmpty() {
				m.quitType = cur.QuitType

				return m, tea.Quit
			}

			cur.IsSelected = !cur.IsSelected

			return m, nil

		case tea.KeyBackspace:
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.updateVisibleOption(false)
			}

		case tea.KeyRunes:
			m.filter = append(m.filter, msg.Runes...)
			m.updateVisibleOption(true)
		}
	}

	// No changes
	return m, nil
}

func (m *Model) updateVisibleOption(add bool) {
	var (
		cur  int
		opts []*innerOption
	)

	if len(m.filter) == 0 {
		opts = slices.Clone(m.innerOpts)
	} else {
		filterStr := strings.ToLower(string(m.filter))

		for _, opt := range helper.If(add, m.innerOptsQueue.Data(), m.innerOpts) {
			if !opt.QuitType.IsEmpty() {
				opts = append(opts, opt)
				continue
			}

			if strings.Contains(strings.ToLower(opt.Text), filterStr) {
				opts = append(opts, opt)
				if opt.IsHovered {
					cur = len(opts) - 1
				}
			}
		}
	}

	m.innerOptsQueue.Current().IsHovered = false
	m.innerOptsQueue = helper.NewCircularQueue(opts, cur)
	m.innerOptsQueue.Current().IsHovered = true
}

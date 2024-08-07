package selectpkg

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	"golang.org/x/term"
	"os"
	"slices"
	"time"
)

type TickMsg struct{}

func (m *Model) Init() tea.Cmd {
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

	m.innerOpts = append(
		helper.SliceMap(m.opts, func(opts Option) *innerOption {
			return &innerOption{
				Text: opts.Text,
				Desc: opts.Desc,
			}
		}),
		&innerOption{Text: m.submitChoiceName, QuitType: QuitTypeSubmit},
		&innerOption{Text: m.quitChoiceName, QuitType: QuitTypeBack},
	)

	m.filter = []rune{}
	m.innerOptsQueue = helper.NewCircularQueue(slices.Clone(m.innerOpts), 0)
	m.innerOptsQueue.Current().IsHovered = true

	m.windowWidth, m.windowHeight, _ = term.GetSize(int(os.Stdin.Fd()))

	tea.ClearScreen()

	m.initPaginator()

	return m.tickCmdFn()
}

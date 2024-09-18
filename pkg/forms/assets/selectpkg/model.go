package selectpkg

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type variant struct {
	Option Option

	IsSelected bool
	IsHovered  bool

	QuitType QuitType
}

type Option struct {
	Text string
	Desc string
}

type Model struct {
	headerFn         func() string
	questionFn       func() string
	quitChoiceName   string
	submitChoiceName string
	reloadInterval   time.Duration
	tickCmdFn        func() tea.Cmd
	opts             []Option

	// Init
	allVariants     []*variant
	filter          []rune
	currentVariants []*variant

	// Result
	quitType QuitType
}

func NewSelectModel(opts []Option) *Model {
	return &Model{
		headerFn:         func() string { return "" },
		questionFn:       func() string { return "" },
		submitChoiceName: DefaultSubmitOption,
		quitChoiceName:   DefaultBackOption,

		opts: opts,
	}
}

func (m *Model) SetHeaderFn(fn func() string) *Model {
	m.headerFn = fn
	return m
}

func (m *Model) SetQuestionFn(fn func() string) *Model {
	m.questionFn = fn
	return m
}

func (m *Model) SetQuitChoiceName(name string) *Model {
	m.quitChoiceName = name
	return m
}

func (m *Model) SetSubmitChoiceName(name string) *Model {
	m.submitChoiceName = name
	return m
}

func (m *Model) SetReloadInterval(interval time.Duration) *Model {
	m.reloadInterval = interval
	return m
}

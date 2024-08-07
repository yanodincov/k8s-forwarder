package selectpkg

import (
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	"time"
)

type innerOption struct {
	Text string
	Desc string

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

	// Static
	tickCmdFn func() tea.Cmd
	opts      []Option

	// Init
	innerOpts      []*innerOption
	filter         []rune
	innerOptsQueue *helper.CircularQueueLimited[*innerOption]
	paginator      paginator.Model
	windowWidth    int
	windowHeight   int

	// Result
	quitType QuitType
}

func NewMultiSelectModel(opts []Option) *Model {
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

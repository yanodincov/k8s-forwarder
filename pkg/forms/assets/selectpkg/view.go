package selectpkg

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	"strings"
	"unicode/utf8"
)

const (
	minPaginatorPageSize = 3
	firstQuitMarinTop    = 1

	blockDelimiterSize = 2
	blocksCount        = 4

	defaultForeground  lipgloss.Color = "#bbbbbb"
	quitForeground     lipgloss.Color = "#777777"
	selectedForeground lipgloss.Color = "#ffffff"

	cursorSymbol          = "→"
	selectedListDelimiter = ", "
)

func (m *Model) View() string {
	var (
		opts         []any
		selectedOpts []string
	)
	for _, opt := range m.innerOptsQueue.Data() {
		if opt.IsSelected {
			selectedOpts = append(selectedOpts, opt.Text)
		}
		opts = append(opts, opt.Text)
	}

	header := m.headerFn()
	question := m.questionFn()
	if len(m.filter) > 0 {
		question += " [" + string(m.filter) + "]"
	}

	if len(selectedOpts) > 0 {
		selectedOptions := "Selected options: "
		identLen := utf8.RuneCountInString(selectedOptions)

		selectedOptions += lipgloss.NewStyle().
			Foreground(selectedForeground).
			Render(helper.GetTextFromItemsWithIdent(helper.MakeColumnTextWithIdentSpec{
				Items:         selectedOpts,
				Delimiter:     selectedListDelimiter,
				MaxRowLen:     m.windowWidth,
				IdentLen:      identLen,
				FirstRowIdent: 0,
			}))

		question += "\n" + selectedOptions
	}

	contentRows := strings.Count(header, "\n") +
		strings.Count(question, "\n") +
		(blocksCount-1)*blockDelimiterSize
	m.rebuildPaginator(helper.Max(m.windowHeight-contentRows, minPaginatorPageSize))
	isSinglePageList := m.innerOptsQueue.Len() <= m.paginator.PerPage

	var iIncrement int
	if !isSinglePageList {
		start, end := m.paginator.GetSliceBounds(m.innerOptsQueue.Len())
		opts = opts[start:end]
		iIncrement = start
	}

	optionList := list.
		New(opts...).
		ItemStyleFunc(func(listItem list.Items, i int) lipgloss.Style {
			i += iIncrement

			style := lipgloss.NewStyle().Foreground(defaultForeground)
			opt := m.innerOptsQueue.Data()[i]

			if opt.IsSelected {
				style = style.Foreground(selectedForeground)
			}
			if opt.IsHovered {
				style = style.Bold(true)
			}
			if !opt.QuitType.IsEmpty() {
				style = style.Italic(true).Foreground(quitForeground)
			} else if isSinglePageList && i < m.innerOptsQueue.Len()-1 && !m.innerOptsQueue.Data()[i+1].QuitType.IsEmpty() {
				style = style.MarginBottom(firstQuitMarinTop)
			}

			return style
		}).
		Enumerator(func(l list.Items, i int) string {
			i += iIncrement

			if m.innerOptsQueue.Data()[i].IsHovered {
				return cursorSymbol
			}

			return " "
		})

	parts := make([]string, blocksCount)
	parts[0] = header
	parts[1] = question
	parts[2] = optionList.String()
	parts[3] = helper.If(isSinglePageList, "", m.paginator.View())

	return strings.Join(parts, strings.Repeat("\n", blockDelimiterSize))
}

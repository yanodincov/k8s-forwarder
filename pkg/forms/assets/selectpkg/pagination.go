package selectpkg

import (
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) initPaginator() {
	m.paginator = paginator.New()
	m.paginator.Type = paginator.Dots
	m.paginator.InactiveDot = lipgloss.NewStyle().Foreground(defaultForeground).Render("■ ")
	m.paginator.ActiveDot = lipgloss.NewStyle().Foreground(selectedForeground).Render("■ ")
}

func (m *Model) rebuildPaginator(pageSize int) {
	m.paginator.PerPage = pageSize
	m.paginator.SetTotalPages(m.innerOptsQueue.Len())
	m.paginator.Page = m.innerOptsQueue.CurI() / m.paginator.PerPage
}

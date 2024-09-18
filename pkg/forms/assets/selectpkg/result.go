package selectpkg

import (
	tea "github.com/charmbracelet/bubbletea"
)

type iResultProvider interface {
	getResult() Result
}

type Result struct {
	SelectedIdx []int
	QuitType    QuitType
}

func (m *Model) getResult() Result {
	var (
		quitType QuitType
		idx      []int
	)
	for i := range m.opts {
		if m.allVariants[i].IsSelected {
			if !m.allVariants[i].QuitType.IsEmpty() {
				quitType = m.allVariants[i].QuitType
				continue
			}

			idx = append(idx, i)
		}
	}

	return Result{
		SelectedIdx: idx,
		QuitType:    quitType,
	}
}

func GetResultFromModel(model tea.Model) Result {
	return model.(iResultProvider).getResult()
}

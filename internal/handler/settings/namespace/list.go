package namespace

import (
	"fmt"
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
)

type ListScreen struct {
	addScreen        *AddScreen
	actionListScreen *ActionListScreen
	repository       *settings.Repository
	errorHandler     *forms.ErrorHandler
}

func NewListScreen(
	addScreen *AddScreen,
	actionListScreen *ActionListScreen,
	repository *settings.Repository,
	errorHandler *forms.ErrorHandler,
) *ListScreen {
	return &ListScreen{
		addScreen:        addScreen,
		actionListScreen: actionListScreen,
		repository:       repository,
		errorHandler:     errorHandler,
	}
}

func (s *ListScreen) Show() error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		settingsModels, err := s.repository.GetSettings()
		if err != nil {
			return nil, err
		}

		items := make([]forms.OptionSpec, 0, len(settingsModels.Namespaces)+2)
		items = append(items, forms.GetCreateOptionSpec("namespace", func(data forms.OptionData) bool {
			if err := s.addScreen.Show(); err != nil {
				s.errorHandler.Handle(err, "failed to show create namespace actions")
			}

			return false
		}))
		for _, ns := range settingsModels.Namespaces {
			ns := ns
			items = append(items, forms.OptionSpec{
				Data: forms.OptionData{
					ID:          ns.ID.String(),
					Name:        ns.Namespace,
					Description: fmt.Sprintf("Config file path: %s", ns.ConfigFilePath),
				},
				Func: func(data forms.OptionData) bool {
					if err := s.actionListScreen.Show(ns.ID); err != nil {
						s.errorHandler.Handle(err, "failed to show namespace actions")
					}

					return false
				},
			})
		}
		items = append(items, forms.CancelOptionSpec)

		return &forms.SelectFormSpec{
			ErrorText:  s.errorHandler.GetErrorText(),
			QuestionFn: "Choose namespace",
			Items:      items,
		}, nil
	})
}

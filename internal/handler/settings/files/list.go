package files

import (
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
)

type ListScreen struct {
	createScreen     *AddScreen
	actionListScreen *ActionListScreen

	repository   *settings.Repository
	errorHandler *forms.ErrorHandler
}

func NewListScreen(
	createScreen *AddScreen,
	actionListScreen *ActionListScreen,
	repository *settings.Repository,
	errorHandler *forms.ErrorHandler,
) *ListScreen {
	return &ListScreen{
		createScreen:     createScreen,
		actionListScreen: actionListScreen,
		repository:       repository,
		errorHandler:     errorHandler,
	}
}

func (s *ListScreen) Show() error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		files, err := s.repository.GetConfigFiles()
		if err != nil {
			return nil, err
		}

		items := make([]forms.OptionSpec, 0, len(files)+2)
		items = append(items, forms.GetCreateOptionSpec("config file path", func(data forms.OptionData) bool {
			if err := s.createScreen.Show(); err != nil {
				s.errorHandler.Handle(err, "Failed to add file path")
				return true
			}

			return false
		}))
		items = append(items, helper.SliceMap(files, func(file settings.ConfigFileSetting) forms.OptionSpec {
			return forms.OptionSpec{
				Data: forms.OptionData{
					ID:          file.ID.String(),
					Name:        file.Path,
					Description: "Manage saved k8s yaml config-file path",
				},
				Func: func(data forms.OptionData) bool {
					if err := s.actionListScreen.Show(file); err != nil {
						s.errorHandler.Handle(err, "Failed to show file actions")
						return true
					}

					return false
				},
			}
		})...)
		items = append(items, forms.CancelOptionSpec)

		return &forms.SelectFormSpec{
			ErrorText:  s.errorHandler.GetErrorText(),
			QuestionFn: "K8s forwarder application",
			Items:      items,
		}, nil
	})
}

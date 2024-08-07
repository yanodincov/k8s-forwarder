package settings

import (
	"github.com/yanodincov/k8s-forwarder/internal/handler/settings/files"
	"github.com/yanodincov/k8s-forwarder/internal/handler/settings/namespace"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
)

type ListScreen struct {
	namespaceListScreen *namespace.ListScreen
	fileListScreen      *files.ListScreen
	errorHandler        *forms.ErrorHandler
}

func NewListScreen(
	namespaceListScreen *namespace.ListScreen,
	fileListScreen *files.ListScreen,
	errorHandler *forms.ErrorHandler,
) *ListScreen {
	return &ListScreen{
		namespaceListScreen: namespaceListScreen,
		fileListScreen:      fileListScreen,
		errorHandler:        errorHandler,
	}
}

func (s *ListScreen) Show() error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		return &forms.SelectFormSpec{
			QuestionFn: "Choose settings",
			ErrorText:  s.errorHandler.GetErrorText(),
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						ID:          "namespaces",
						Name:        "Saved k8s namespaces",
						Description: "Manage saved k8s namespaces",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.namespaceListScreen.Show(); err != nil {
							s.errorHandler.Handle(err, "failed to show namespaces settings")
						}

						return false
					},
				},
				{
					Data: forms.OptionData{
						ID:          "files",
						Name:        "Saved k8s yaml config-files",
						Description: "Manage k8s config files",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.fileListScreen.Show(); err != nil {
							s.errorHandler.Handle(err, "failed to show files settings")
						}

						return false
					},
				},
				forms.CancelOptionSpec,
			},
		}, nil
	})
}

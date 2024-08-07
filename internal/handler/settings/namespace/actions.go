package namespace

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
)

type ActionListScreen struct {
	deleteScreen *DeleteScreen
	repository   *settings.Repository
	errorHandler *forms.ErrorHandler
}

func NewActionListScreen(
	deleteScreen *DeleteScreen,
	repository *settings.Repository,
	errorHandler *forms.ErrorHandler,
) *ActionListScreen {
	return &ActionListScreen{
		deleteScreen: deleteScreen,
		repository:   repository,
		errorHandler: errorHandler,
	}
}

func (s *ActionListScreen) Show(id uuid.UUID) error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		namespace, err := s.repository.GetNamespace(id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get namespace")
		}

		return &forms.SelectFormSpec{
			HeaderFn: forms.GetTableInfoWithTitle("Manage namespace",
				"Namespace", namespace.Namespace,
				"Config path", namespace.ConfigFilePath,
			),
			ErrorText:  s.errorHandler.GetErrorText(),
			QuestionFn: "Select action",
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						ID:          "delete",
						Name:        "Delete",
						Description: "Delete namespace",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.deleteScreen.Show(namespace); err != nil {
							s.errorHandler.Handle(err, "failed delete namespace form")
							return false
						}

						return true
					},
				},
				forms.CancelOptionSpec,
			},
		}, nil
	})
}

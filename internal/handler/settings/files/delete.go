package files

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
)

type DeleteScreen struct {
	repository   *settings.Repository
	errorHandler *forms.ErrorHandler
}

func NewDeleteScreen(
	repository *settings.Repository,
	errorHandler *forms.ErrorHandler,
) *DeleteScreen {
	return &DeleteScreen{
		repository:   repository,
		errorHandler: errorHandler,
	}
}

func (s *DeleteScreen) Show(id uuid.UUID) (bool, error) {
	isDeleted := false

	err := forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		configFile, err := s.repository.GetConfigFile(id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get config file")
		}

		return &forms.SelectFormSpec{
			QuestionFn: "Do you want to delete file '" + configFile.Path + "'?",
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						ID:   "delete",
						Name: "Delete",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.repository.RemoveConfigFile(id); err != nil {
							s.errorHandler.Handle(err, "failed to remove file")
							return true
						}
						isDeleted = true

						return true
					},
				},
				forms.CancelOptionSpec,
			},
		}, nil
	})

	return isDeleted, err
}

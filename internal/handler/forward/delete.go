package forward

import (
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
)

type DeleteScreen struct {
	repository   *portset.Repository
	errorHandler *forms.ErrorHandler
}

func NewDeleteScreen(
	repository *portset.Repository,
	errorHandler *forms.ErrorHandler,
) *DeleteScreen {
	return &DeleteScreen{
		repository:   repository,
		errorHandler: errorHandler,
	}
}

func (s *DeleteScreen) Show(config *portset.ServiceSet) error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		return &forms.SelectFormSpec{
			QuestionFn: "Do you want to delete port set '" + config.Name + "'?",
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						ID:          "yes",
						Name:        "Delete",
						Description: "Delete port set",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.repository.RemoveSet(config.ID); err != nil {
							s.errorHandler.Handle(err, "failed to delete port set")
							return true
						}

						return true
					},
				},
				forms.CancelOptionSpec,
			},
		}, nil
	})
}

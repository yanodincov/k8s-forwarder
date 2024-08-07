package namespace

import (
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
)

type DeleteScreen struct {
	repository   *settings.Repository
	errorHandler *forms.ErrorHandler
}

func NewDeleteScreen(repository *settings.Repository, errorHandler *forms.ErrorHandler) *DeleteScreen {
	return &DeleteScreen{repository: repository, errorHandler: errorHandler}
}

func (s *DeleteScreen) Show(namespace *settings.NamespaceSetting) error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		return &forms.SelectFormSpec{
			QuestionFn: "Do you want to delete namespace '" + namespace.Namespace + "'?",
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						ID:          "yes",
						Name:        "Delete",
						Description: "Delete namespace",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.repository.RemoveNamespace(namespace.ID); err != nil {
							s.errorHandler.Handle(err, "failed to delete namespace")
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

package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
)

type ListScreen struct {
	addScreen    *AddScreen
	actionScreen *ActionsScreen
	repository   *portset.Repository
	errorHandler *forms.ErrorHandler
}

func NewListScreen(
	addScreen *AddScreen,
	actionScreen *ActionsScreen,
	repository *portset.Repository,
	errorHandler *forms.ErrorHandler,
) *ListScreen {
	return &ListScreen{
		addScreen:    addScreen,
		actionScreen: actionScreen,
		repository:   repository,
		errorHandler: errorHandler,
	}
}

func (s *ListScreen) Show(setID uuid.UUID) error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		set, err := s.repository.GetServiceSet(setID)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get set list")
		}
		if set == nil {
			return nil, errors.New("port set not found")
		}

		options := make([]forms.OptionSpec, 0, len(set.Services)+2)
		options = append(options, forms.GetCreateOptionSpec("service", func(data forms.OptionData) bool {
			if err := s.addScreen.Show(setID); err != nil {
				s.errorHandler.Handle(err)
			}

			return false
		}))
		options = append(options, helper.SliceMap(set.Services, func(service portset.ServiceForwardConfig) forms.OptionSpec {
			return forms.OptionSpec{
				Data: forms.OptionData{
					ID: service.ID.String(),
					Name: fmt.Sprintf("%s/%s:%d <-> %d",
						service.Namespace, service.ServiceName, service.ServicePort, service.LocalPort),
				},
				Func: func(data forms.OptionData) bool {
					if err := s.actionScreen.Show(&service); err != nil {
						s.errorHandler.Handle(err)
					}

					return false
				},
			}
		})...)
		options = append(options, forms.CancelOptionSpec)

		return &forms.SelectFormSpec{
			ErrorText:  s.errorHandler.GetErrorText(),
			QuestionFn: "Manage port set '" + set.Name + "' services",
			Items:      options,
		}, nil
	})
}

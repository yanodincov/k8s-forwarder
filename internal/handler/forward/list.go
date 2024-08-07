package forward

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/internal/service/forwarder"
	"github.com/yanodincov/k8s-forwarder/pkg/cli"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	"golang.org/x/exp/maps"
)

type ListScreen struct {
	actionsScreen *ActionsScreen
	createScreen  *CreateScreen
	deleteScreen  *DeleteScreen

	forwardService *forwarder.Service

	repository   *portset.Repository
	errorHandler *forms.ErrorHandler
}

func NewListScreen(
	actionsScreen *ActionsScreen,
	createScreen *CreateScreen,
	deleteScreen *DeleteScreen,
	forwardService *forwarder.Service,
	repository *portset.Repository,
	errorHandler *forms.ErrorHandler,
) *ListScreen {
	return &ListScreen{
		actionsScreen:  actionsScreen,
		createScreen:   createScreen,
		deleteScreen:   deleteScreen,
		forwardService: forwardService,
		repository:     repository,
		errorHandler:   errorHandler,
	}
}

func (s *ListScreen) Show() error {
	cli.ClearScreen()

	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		sets, err := s.repository.GetServiceSets()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get set list")
		}

		setByIDs := helper.Slice2Map(sets.Sets, func(set portset.ServiceSet) (uuid.UUID, portset.ServiceSet) {
			return set.ID, set
		})

		var header []string
		for setID, status := range s.forwardService.GetServiceSetsForwardStatus(maps.Keys(setByIDs)) {
			set := setByIDs[setID]
			header = append(header, set.Name, status.StringWithColor())
		}

		list := make([]forms.OptionSpec, 0, len(sets.Sets)+2)
		list = append(list, forms.GetCreateOptionSpec("port set", func(data forms.OptionData) bool {
			if showErr := s.createScreen.Show(); showErr != nil {
				s.errorHandler.Handle(showErr)
			}

			return false
		}))

		for _, set := range sets.Sets {
			config := set
			list = append(list, forms.OptionSpec{
				Data: forms.OptionData{
					ID:   config.ID.String(),
					Name: config.Name,
				},
				Func: func(data forms.OptionData) bool {
					if showErr := s.actionsScreen.Show(&config); showErr != nil {
						s.errorHandler.Handle(showErr)
					}

					return false
				},
			})
		}

		list = append(list, forms.CancelOptionSpec)

		return &forms.SelectFormSpec{
			HeaderFn:   forms.GetTableInfoWithTitle("Port sets", header...),
			QuestionFn: "Choose configuration",
			ErrorText:  s.errorHandler.GetErrorText(),
			Items:      list,
		}, nil
	})
}

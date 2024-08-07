package forward

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/yanodincov/k8s-forwarder/internal/handler/forward/services"
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/internal/service/forwarder"
	"github.com/yanodincov/k8s-forwarder/pkg/cli"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
)

type ActionsScreen struct {
	serviceListScreen *services.ListScreen
	deleteScreen      *DeleteScreen
	forwardService    *forwarder.Service
	repository        *portset.Repository
	errorHandler      *forms.ErrorHandler
}

func NewActionsScreen(
	serviceListScreen *services.ListScreen,
	deleteScreen *DeleteScreen,
	forwardService *forwarder.Service,
	repository *portset.Repository,
	errorHandler *forms.ErrorHandler,
) *ActionsScreen {
	return &ActionsScreen{
		serviceListScreen: serviceListScreen,
		deleteScreen:      deleteScreen,
		forwardService:    forwardService,
		repository:        repository,
		errorHandler:      errorHandler,
	}
}

func (s *ActionsScreen) Show(config *portset.ServiceSet) error {
	cli.ClearScreen()

	statusByServiceID := s.forwardService.GetServicesForwardStatus(
		helper.SliceMap(config.Services, func(service portset.ServiceForwardConfig) uuid.UUID {
			return service.ID
		}),
	)

	headers := make([]string, 0, len(config.Services)*2)
	for _, service := range config.Services {
		serviceName := fmt.Sprintf("%s/%s:%d <-> %d",
			service.Namespace, service.ServiceName, service.ServicePort, service.LocalPort)
		headers = append(headers, serviceName, statusByServiceID[service.ID].StringWithColor())
	}

	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		return &forms.SelectFormSpec{
			HeaderFn:   forms.GetTableInfoWithTitle("Manage "+config.Name+" port set", headers...),
			QuestionFn: "Choose action",
			ErrorText:  s.errorHandler.GetErrorText(),
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						Name:        "Start forward",
						Description: "Forward all set services",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.forwardService.ForwardSet(config.ID); err != nil {
							s.errorHandler.Handle(err, "failed to start forward service")
						}

						return true
					},
				},
				{
					Data: forms.OptionData{
						Name:        "Stop forward",
						Description: "Stop all set services",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.forwardService.StopForwardSet(config.ID); err != nil {
							s.errorHandler.Handle(err, "failed to stop forward service")
						}

						return true
					},
				},
				{
					Data: forms.OptionData{
						Name:        "Services",
						Description: "Manage port set services",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.serviceListScreen.Show(config.ID); err != nil {
							s.errorHandler.Handle(err, "failed to show service list")
						}

						return true
					},
				},
				{
					Data: forms.OptionData{
						ID:          "delete",
						Name:        "Delete",
						Description: "Delete set",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.deleteScreen.Show(config); err != nil {
							s.errorHandler.Handle(err, "failed to show delete form")
						}

						return true
					},
				},
				forms.CancelOptionSpec,
			},
		}, nil
	})
}

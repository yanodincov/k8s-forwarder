package services

import (
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/internal/service/forwarder"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"strconv"
)

type DeleteScreen struct {
	forwardService *forwarder.Service
	repository     *portset.Repository
	errorHandler   *forms.ErrorHandler
}

func NewDeleteScreen(
	forwardService *forwarder.Service,
	repository *portset.Repository,
	errorHandler *forms.ErrorHandler,
) *DeleteScreen {
	return &DeleteScreen{
		forwardService: forwardService,
		repository:     repository,
		errorHandler:   errorHandler,
	}
}
func (s *DeleteScreen) Show(service *portset.ServiceForwardConfig) error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		return &forms.SelectFormSpec{
			HeaderFn: forms.GetTableInfoWithTitle("Service",
				"Config path", service.ConfigFilePath,
				"Namespace", service.Namespace,
				"Service", service.ServiceName,
				"Service port", strconv.Itoa(service.ServicePort),
				"Local port", strconv.Itoa(service.LocalPort),
			),
			QuestionFn: "Do you want to delete service '" + service.ServiceName + "'?",
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						ID:          "delete",
						Name:        "Delete",
						Description: "Remove service from port set",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.repository.RemoveService(service.ID); err != nil {
							s.errorHandler.Handle(err, "failed to remove service")
							return true
						}
						s.forwardService.StopForwardService(service.ID)

						return true
					},
				},
				forms.CancelOptionSpec,
			},
		}, nil
	})
}

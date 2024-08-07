package services

import (
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"strconv"
)

type ActionsScreen struct {
	deleteScreen *DeleteScreen
	errorHandler *forms.ErrorHandler
}

func NewActionsScreen(deleteScreen *DeleteScreen, errorHandler *forms.ErrorHandler) *ActionsScreen {
	return &ActionsScreen{deleteScreen: deleteScreen, errorHandler: errorHandler}
}

func (s *ActionsScreen) Show(service *portset.ServiceForwardConfig) error {

	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		return &forms.SelectFormSpec{
			HeaderFn: forms.GetTableInfoWithTitle("Service",
				"Config path", service.ConfigFilePath,
				"Namespace", service.Namespace,
				"Service", service.ServiceName,
				"Service port", strconv.Itoa(service.ServicePort),
				"Local port", strconv.Itoa(service.LocalPort),
			),
			QuestionFn: "Select action",
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						ID:          "delete",
						Name:        "Delete",
						Description: "Remove service from port set",
					},
					Func: func(data forms.OptionData) bool {
						if err := s.deleteScreen.Show(service); err != nil {
							s.errorHandler.Handle(err, "failed to remove service")
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

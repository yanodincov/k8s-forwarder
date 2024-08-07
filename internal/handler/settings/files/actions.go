package files

import (
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
)

type ActionListScreen struct {
	deleteScreen *DeleteScreen
	errorHandler *forms.ErrorHandler
}

func NewActionListScreen(
	deleteScreen *DeleteScreen,
	errorHandler *forms.ErrorHandler,
) *ActionListScreen {
	return &ActionListScreen{
		deleteScreen: deleteScreen,
		errorHandler: errorHandler,
	}
}

func (s *ActionListScreen) Show(file settings.ConfigFileSetting) error {
	return forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
		return &forms.SelectFormSpec{
			HeaderFn: forms.GetTableInfoWithTitle("Manage saved k8s config path",
				"Config path", file.Path,
			),
			QuestionFn: "Select action",
			ErrorText:  s.errorHandler.GetErrorText(),
			Items: []forms.OptionSpec{
				{
					Data: forms.OptionData{
						ID:          "delete",
						Name:        "Delete",
						Description: "Remove file from saved k8s-forwarder k8s config files (this not delete file from disk)",
					},
					Func: func(data forms.OptionData) bool {
						isDeleted, err := s.deleteScreen.Show(file.ID)
						if err != nil {
							s.errorHandler.Handle(err, "Failed to delete file")
							return false
						}

						return isDeleted
					},
				},
				forms.CancelOptionSpec,
			},
		}, nil
	})
}

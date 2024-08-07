package files

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	k8sclient "github.com/yanodincov/k8s-forwarder/pkg/k8s-client"
	"os/user"
	"path/filepath"
	"strings"
)

type AddScreen struct {
	repository   *settings.Repository
	errorHandler *forms.ErrorHandler
}

func NewAddScreen(
	repository *settings.Repository,
	errorHandler *forms.ErrorHandler,
) *AddScreen {
	return &AddScreen{
		repository:   repository,
		errorHandler: errorHandler,
	}
}

func (s *AddScreen) Show() error {
	questions := []*survey.Question{
		{
			Name: "path",
			Prompt: &survey.Input{
				Message: "Enter k8s config file path:",
				Suggest: func(toComplete string) []string {
					// Expand ~ to the home directory
					if strings.HasPrefix(toComplete, "~") {
						if usr, err := user.Current(); err == nil {
							toComplete = filepath.Join(usr.HomeDir, strings.TrimPrefix(toComplete, "~"))
						}
					}
					matches, _ := filepath.Glob(toComplete + "*")

					return matches
				},
			},
			Validate: func(val interface{}) error {
				path, ok := val.(string)
				if !ok || path == "" {
					return errors.New("file path is required")
				}

				if err := k8sclient.ValidateK8sConfigFile(path); err != nil {
					return errors.Wrap(err, "validate k8s config file")
				}

				existFile, err := s.repository.GetConfigFileByPath(path)
				if err != nil {
					return errors.Wrap(err, "check file exist")
				}

				if existFile != nil {
					return errors.New("file already exist")
				}

				return nil
			},
		},
	}

	file := settings.ConfigFileSetting{
		ID: uuid.New(),
	}
	if err := survey.Ask(questions, &file); err != nil {
		if forms.IsInterruptSurveyErr(err) {
			return nil
		}
		return errors.Wrap(err, "configure k8s config file")
	}

	if err := s.repository.AddConfigFile(file); err != nil {
		return errors.Wrap(err, "save k8s config file")
	}

	return nil
}

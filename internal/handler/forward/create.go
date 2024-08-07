package forward

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/internal/service/k8s"
	"github.com/yanodincov/k8s-forwarder/pkg/cli"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
)

type CreateScreen struct {
	k8sService   *k8s.Service
	repository   *portset.Repository
	errorHandler *forms.ErrorHandler
}

func NewCreateScreen(
	k8sService *k8s.Service,
	repository *portset.Repository,
	errorHandler *forms.ErrorHandler,
) *CreateScreen {
	return &CreateScreen{
		k8sService:   k8sService,
		repository:   repository,
		errorHandler: errorHandler,
	}
}

func (c *CreateScreen) Show() error {
	cli.ClearScreen()

	set := portset.ServiceSet{
		ID: uuid.New(),
	}

	questions := []*survey.Question{
		{
			Name: "Name",
			Prompt: &survey.Input{
				Message: "Enter port set name: ",
			},
			Validate: func(val interface{}) error {
				str, ok := val.(string)
				if !ok {
					return errors.New("invalid input value type")
				}

				foundedSet, err := c.repository.GetServiceSetByName(str)
				if err != nil {
					return errors.Wrap(err, "failed to get port set by name")
				}

				if foundedSet != nil {
					return errors.New("port set with this name already exists")
				}

				return nil
			},
		},
	}

	if err := survey.Ask(questions, &set, forms.GetSurveyOpts()...); err != nil {
		if forms.IsInterruptSurveyErr(err) {
			return nil
		}
		return errors.Wrap(err, "failed to ask questions")
	}

	if err := c.repository.AddSet(&set); err != nil {
		return errors.Wrap(err, "failed to create port set")
	}

	return nil
}

package namespace

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/internal/service/k8s"
	"github.com/yanodincov/k8s-forwarder/pkg/cli"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	k8sclient "github.com/yanodincov/k8s-forwarder/pkg/k8s-client"
	"strings"
)

const (
	k8sLackPermissionText = "cannot list resource \"namespaces\""
)

type AddScreen struct {
	service            *k8s.Service
	settingsRepository *settings.Repository
	errorHandler       *forms.ErrorHandler
}

func NewAddScreen(
	service *k8s.Service,
	settingsRepository *settings.Repository,
	errorHandler *forms.ErrorHandler,
) *AddScreen {
	return &AddScreen{
		service:            service,
		settingsRepository: settingsRepository,
		errorHandler:       errorHandler,
	}
}

type addNamespacesSpec struct {
	ConfigFilePath string
	Namespace      string
}

func (s *AddScreen) Show() error {
	cli.ClearScreen()
	fmt.Println("")

	configFiles, err := s.service.GetConfigFilePathsWithDefault()
	if err != nil {
		return errors.Wrap(err, "failed to get config files")
	}
	if len(configFiles) == 0 {
		return errors.New("config file path not found: add k8s yaml config file path")
	}

	questions := []*survey.Question{
		{
			Name: "ConfigFilePath",
			Prompt: &survey.Select{
				Message: "Select k8s yaml config file:",
				Options: configFiles,
			},
			Validate: func(opt interface{}) error {
				var err error

				answer, ok := opt.(core.OptionAnswer)
				if !ok {
					return errors.New("invalid option")
				}
				path := answer.Value

				if err = k8sclient.ValidateK8sConfigFile(path); err != nil {
					return err
				}

				return nil
			},
		},
	}

	askOpts := forms.GetSurveyOpts()

	spec := addNamespacesSpec{}
	if err = survey.Ask(questions, &spec, askOpts...); err != nil {
		if forms.IsInterruptSurveyErr(err) {
			return nil
		}
		return errors.Wrap(err, "failed to ask questions")
	}

	namespaces, err := s.service.GetNamespaces(spec.ConfigFilePath)

	errText := helper.IfFnOrDef(err != nil, func() string {
		return err.Error()
	})
	if strings.Contains(errText, k8sLackPermissionText) {
		errText = "lack of permission to list namespaces"
	} else if len(namespaces) == 0 {
		errText = "not found namespaces via selected k8s config file"
	}

	if errText != "" {
		fmt.Println(" Discover namespaces form k8s error: " + promptui.Styler(promptui.FGRed)(errText) +
			"\n " + promptui.Styler(promptui.FGYellow)("Suggestion disabled"))
	}

	if err = survey.Ask([]*survey.Question{
		{
			Name: "Namespace",
			Prompt: &survey.Input{
				Message: "Enter namespaces:",
				Suggest: helper.IfFnOrDef(len(namespaces) > 0, func() func(toComplete string) []string {
					return func(toComplete string) []string {
						return helper.SliceFilter(namespaces, func(input string) bool {
							return strings.Contains(input, toComplete)
						})
					}
				}),
			},
			Validate: func(val interface{}) error {
				namespace, ok := val.(string)
				if !ok || namespace == "" {
					return errors.New("namespace is required")
				}

				_, err := s.service.GetServices(spec.ConfigFilePath, namespace)
				if err != nil {
					return errors.Wrap(err, "failed to get k8s services")
				}

				existedNamespace, err := s.settingsRepository.GetNamespaceByNameAndFile(namespace, spec.ConfigFilePath)
				if err != nil {
					return errors.Wrap(err, "failed to get namespace")
				}
				if existedNamespace != nil {
					return errors.New("namespace already exist")
				}

				return nil
			},
		},
	}, &spec, askOpts...); err != nil {
		if forms.IsInterruptSurveyErr(err) {
			return nil
		}
		return errors.Wrap(err, "failed to ask questions")
	}

	if err = s.settingsRepository.AddNamespace(settings.NamespaceSetting{
		ID:             uuid.New(),
		Namespace:      spec.Namespace,
		ConfigFilePath: spec.ConfigFilePath,
	}); err != nil {
		return errors.Wrap(err, "failed to create namespace")
	}

	return nil
}

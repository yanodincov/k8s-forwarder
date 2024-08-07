package services

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/internal/service/k8s"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	v1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

type AddScreen struct {
	k8sService         *k8s.Service
	repository         *portset.Repository
	settingsRepository *settings.Repository
}

func NewAddScreen(
	k8sService *k8s.Service,
	repository *portset.Repository,
	settingsRepository *settings.Repository,
) *AddScreen {
	return &AddScreen{
		k8sService:         k8sService,
		repository:         repository,
		settingsRepository: settingsRepository,
	}
}

type servicePortSpec struct {
	ServicePort string
}

func (s *AddScreen) Show(setID uuid.UUID) error {
	namespaces, err := s.settingsRepository.GetNamespaces()
	if err != nil {
		return errors.Wrap(err, "get namespaces")
	}

	questions := []*survey.Question{
		{
			Name: "Namespace",
			Prompt: &survey.Select{
				Message: "Select namespace:",
				Options: helper.SliceMap(namespaces, func(namespace settings.NamespaceSetting) string {
					return namespace.Namespace
				}),
			},
			Validate: survey.Required,
		},
	}

	serviceForwardCfg := portset.ServiceForwardConfig{
		ID:    uuid.New(),
		SetID: setID,
	}
	if err := survey.Ask(questions, &serviceForwardCfg); err != nil {
		if forms.IsInterruptSurveyErr(err) {
			return nil
		}
		return errors.Wrap(err, "configure k8s config service")
	}

	selectedNamespace, ok := helper.SliceFind(namespaces, func(namespace settings.NamespaceSetting) bool {
		return namespace.Namespace == serviceForwardCfg.Namespace
	})
	if !ok {
		return errors.New("namespace not found")
	}
	serviceForwardCfg.ConfigFilePath = selectedNamespace.ConfigFilePath

	k8sServices, err := s.k8sService.GetServices(selectedNamespace.ConfigFilePath, selectedNamespace.Namespace)
	if err != nil {
		return errors.Wrap(err, "get k8s service")
	}

	questions = []*survey.Question{
		{
			Name: "ServiceName",
			Prompt: &survey.Select{
				Message: "Select service:",
				Options: helper.SliceMap(k8sServices, func(service v1.Service) string {
					return service.Name
				}),
			},
		},
	}
	if err := survey.Ask(questions, &serviceForwardCfg); err != nil {
		if forms.IsInterruptSurveyErr(err) {
			return nil
		}
		return errors.Wrap(err, "configure k8s service")
	}

	k8sService, ok := helper.SliceFind(k8sServices, func(service v1.Service) bool {
		return service.Name == serviceForwardCfg.ServiceName
	})
	if !ok {
		return errors.New("service not found")
	}

	questions = []*survey.Question{
		{
			Name: "ServicePort",
			Prompt: &survey.Select{
				Message: "Select service port:",
				Options: helper.SliceMap(k8sService.Spec.Ports, func(port v1.ServicePort) string {
					return port.Name + " : " + strconv.Itoa(int(port.Port))
				}),
			},
			Validate: func(val interface{}) error {
				option, ok := val.(survey.OptionAnswer)
				if !ok {
					return errors.New("service port is required")
				}

				optionParts := strings.Split(option.Value, " : ")
				if len(optionParts) != 2 {
					return errors.New("invalid service port format")
				}

				port, err := strconv.Atoi(optionParts[1])
				if err != nil {
					return errors.Wrap(err, "parse service port")
				}

				service, err := s.repository.GetOneServicePortByFilter(portset.ServicePortFilter{
					SetID:       helper.Ptr(setID),
					ConfigPath:  helper.Ptr(serviceForwardCfg.ConfigFilePath),
					Namespace:   helper.Ptr(serviceForwardCfg.Namespace),
					ServiceName: helper.Ptr(serviceForwardCfg.ServiceName),
					ServicePort: helper.Ptr(port),
				})
				if err != nil {
					return errors.Wrap(err, "get service port")
				}
				if service != nil {
					return errors.New("service port already used in current set")
				}

				return nil
			},
			Transform: func(val interface{}) interface{} {
				option := val.(survey.OptionAnswer)

				return survey.OptionAnswer{
					Value: strings.Split(option.Value, " : ")[1],
					Index: option.Index,
				}
			},
		},
	}

	var portSpec servicePortSpec
	if err := survey.Ask(questions, &portSpec); err != nil {
		if forms.IsInterruptSurveyErr(err) {
			return nil
		}
		return errors.Wrap(err, "configure k8s service port")
	}

	serviceForwardCfg.ServicePort, err = strconv.Atoi(portSpec.ServicePort)
	if err != nil {
		return errors.Wrap(err, "parse service port")
	}

	questions = []*survey.Question{
		{
			Name: "LocalPort",
			Prompt: &survey.Input{
				Message: "Enter local port:",
			},
			Validate: func(val interface{}) error {
				port, ok := val.(string)
				if !ok {
					return errors.New("local port is required")
				}

				_, err := strconv.Atoi(port)
				if err != nil {
					return errors.Wrap(err, "local port is not a number")
				}

				service, err := s.repository.GetOneServicePortByFilter(portset.ServicePortFilter{
					SetID:     helper.Ptr(setID),
					LocalPort: helper.Ptr(helper.Must(strconv.Atoi(port))),
				})
				if err != nil {
					return errors.Wrap(err, "get local port")
				}
				if service != nil {
					return errors.New("local port already exist in current set")
				}

				return nil
			},
		},
	}

	if err := survey.Ask(questions, &serviceForwardCfg); err != nil {
		if forms.IsInterruptSurveyErr(err) {
			return nil
		}
		return errors.Wrap(err, "configure local port")
	}

	if err := s.repository.AddService(serviceForwardCfg); err != nil {
		return errors.Wrap(err, "save service")
	}

	return nil
}

package forwarder

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/internal/service/k8s"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"math/rand"
	"net/http"
	"time"
)

type ForwardStatus int

func (s ForwardStatus) String() string {
	switch s {
	case ForwardStatusActive:
		return "Active"
	case ForwardStatusConnecting:
		return "Connecting"
	case ForwardStatusInactive:
		return "Inactive"
	case ForwardStatusFailed:
		return "Failed"
	default:
		return "Unknown"
	}
}

func (s ForwardStatus) StringWithColor() string {
	switch s {
	case ForwardStatusActive:
		return promptui.Styler(promptui.FGGreen)("Active")
	case ForwardStatusConnecting:
		return promptui.Styler(promptui.FGYellow)("Connecting")
	case ForwardStatusInactive:
		return promptui.Styler(promptui.FGBlue)("Inactive")
	case ForwardStatusFailed:
		return promptui.Styler(promptui.FGRed)("Failed")
	default:
		return "Unknown"
	}
}

const (
	ForwardStatusActive ForwardStatus = iota
	ForwardStatusConnecting
	ForwardStatusInactive
	ForwardStatusFailed
)

type forwardState struct {
	Status ForwardStatus
	Stop   context.CancelFunc
}

func newForwardState(status ForwardStatus, stop context.CancelFunc) forwardState {
	return forwardState{Status: status, Stop: stop}
}

type Service struct {
	forwardedPorts *helper.SyncMap[uuid.UUID, forwardState]
	k8sService     *k8s.Service
	repository     *portset.Repository
}

func NewService(repository *portset.Repository, k8sService *k8s.Service) *Service {
	return &Service{
		forwardedPorts: helper.NewSyncMap[uuid.UUID, forwardState](),
		k8sService:     k8sService,
		repository:     repository,
	}
}

type clientConfig struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

const forwardSetTiemout = 15 * time.Second

func (s *Service) ForwardSet(setID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), forwardSetTiemout)
	defer cancel()

	set, err := s.repository.GetServiceSet(setID)
	if err != nil {
		return errors.Wrap(err, "get service set")
	}

	clientsByConfigFile := make(map[string]clientConfig)
	for _, service := range set.Services {
		clientWithConfig, ok := clientsByConfigFile[service.ConfigFilePath]
		if !ok {
			config, err := clientcmd.BuildConfigFromFlags("", service.ConfigFilePath)
			if err != nil {
				return errors.Wrap(err, "build kubeconfig")
			}

			client, err := kubernetes.NewForConfig(config)
			if err != nil {
				return errors.Wrap(err, "create clientset")
			}

			clientWithConfig = clientConfig{clientset: client, config: config}
			clientsByConfigFile[service.ConfigFilePath] = clientWithConfig
		}

		if err = s.portForward(ctx, clientWithConfig, service); err != nil {
			return errors.Wrap(err, "port forward")
		}
	}

	return nil
}

func (s *Service) getPodName(ctx context.Context, clientCfg clientConfig, namespace string, serviceName string) (string, error) {
	service, err := clientCfg.clientset.CoreV1().Services(namespace).Get(ctx, serviceName, v1.GetOptions{})
	if err != nil {
		return "", errors.Wrap(err, "get service")
	}

	selector := service.Spec.Selector

	selectorString := ""
	for key, value := range selector {
		if selectorString != "" {
			selectorString += ","
		}
		selectorString += fmt.Sprintf("%s=%s", key, value)
	}

	pods, err := clientCfg.clientset.CoreV1().Pods(namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selectorString,
	})
	if err != nil {
		return "", errors.Wrap(err, "list pods")
	}

	if len(pods.Items) == 0 {
		return "", errors.New("no pods found")
	}

	rand.Seed(time.Now().UnixNano())
	randomPod := pods.Items[rand.Intn(len(pods.Items))]

	return randomPod.Name, nil
}

func (s *Service) portForward(
	ctx context.Context,
	clientCfg clientConfig,
	service portset.ServiceForwardConfig,
) error {
	// Добавить проверку на то, что форвардинг уже запущен

	// Получаем имя пода по имени сервиса
	podName, err := s.getPodName(ctx, clientCfg, service.Namespace, service.ServiceName)
	if err != nil {
		return errors.Wrap(err, "get pod name")
	}

	// Получаем RESTClient
	restClient := clientCfg.clientset.CoreV1().RESTClient()

	// Получаем URL для SPDY транспорта
	url := restClient.Post().
		Resource("pods").
		Namespace(service.Namespace).
		Name(podName).
		SubResource("portforward").
		URL()

	// Настраиваем SPDY транспорт
	transport, upgrader, err := spdy.RoundTripperFor(clientCfg.config)
	if err != nil {
		return errors.Wrap(err, "create round tripper")
	}

	// Настраиваем порты для проброса
	ports := []string{fmt.Sprintf("%d:%d", service.LocalPort, service.ServicePort)}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.forwardedPorts.Set(service.ID, newForwardState(ForwardStatusConnecting, cancel))

	readyCh := make(chan struct{})
	go func() {
		port, _ := s.forwardedPorts.Get(service.ID)

		select {
		case <-ctx.Done():
			port.Status = ForwardStatusInactive
		case <-time.After(forwardSetTiemout):
			port.Status = ForwardStatusFailed
			cancel()
		case <-readyCh:
			port.Status = ForwardStatusActive
		}

		s.forwardedPorts.Set(service.ID, port)
	}()

	// Настраиваем portforward
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, url)
	fw, err := portforward.New(dialer, ports, ctx.Done(), readyCh, io.Discard, io.Discard)
	if err != nil {
		return errors.Wrap(err, "create portforward")
	}

	// Запускаем portforward
	go func() {
		port, _ := s.forwardedPorts.Get(service.ID)

		if err = fw.ForwardPorts(); err != nil {
			port.Status = ForwardStatusFailed
		} else {
			port.Status = ForwardStatusActive
		}

		s.forwardedPorts.Set(service.ID, port)
	}()

	return nil
}

func (s *Service) GetServicesForwardStatus(ids []uuid.UUID) map[uuid.UUID]ForwardStatus {
	result := make(map[uuid.UUID]ForwardStatus, len(ids))
	for _, id := range ids {
		status, ok := s.forwardedPorts.Get(id)
		if !ok {
			result[id] = ForwardStatusInactive
		} else {
			result[id] = status.Status
		}
	}

	return result
}

func (s *Service) GetServiceSetsForwardStatus(setIDs []uuid.UUID) map[uuid.UUID]ForwardStatus {
	return helper.Slice2Map(setIDs, func(setID uuid.UUID) (uuid.UUID, ForwardStatus) {
		return setID, s.GetServiceSetForwardStatus(setID)
	})
}

func (s *Service) GetServiceSetForwardStatus(setID uuid.UUID) ForwardStatus {
	set, err := s.repository.GetServiceSet(setID)
	if err != nil {
		return ForwardStatusFailed
	}

	for _, service := range set.Services {
		status, ok := s.forwardedPorts.Get(service.ID)
		if !ok {
			return ForwardStatusInactive
		}

		if status.Status == ForwardStatusFailed {
			return ForwardStatusFailed
		}
	}

	return ForwardStatusActive
}

func (s *Service) StopForwardSet(setID uuid.UUID) error {
	set, err := s.repository.GetServiceSet(setID)
	if err != nil {
		return errors.Wrap(err, "get service set")
	}

	for _, service := range set.Services {
		state, ok := s.forwardedPorts.Get(service.ID)
		if !ok {
			continue
		}

		state.Stop()
		s.forwardedPorts.Delete(service.ID)
	}

	return nil
}

func (s *Service) StopForwardService(serviceID uuid.UUID) {
	state, ok := s.forwardedPorts.Get(serviceID)
	if !ok {
		return
	}

	state.Stop()
	s.forwardedPorts.Delete(serviceID)

	return
}

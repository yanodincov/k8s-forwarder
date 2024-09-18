package handler

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yanodincov/k8s-forwarder/internal/handler/forward"
	"github.com/yanodincov/k8s-forwarder/internal/handler/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/cli"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"github.com/yanodincov/k8s-forwarder/pkg/forms/assets/selectpkg"
	"time"
)

type MainScreen struct {
	errorHandler    *forms.ErrorHandler
	portSetHandler  *forward.ListScreen
	settingsHandler *settings.ListScreen
}

func NewMainScreen(
	errorHandler *forms.ErrorHandler,
	portSetHandler *forward.ListScreen,
	settingsHandler *settings.ListScreen,
) *MainScreen {
	return &MainScreen{
		errorHandler:    errorHandler,
		portSetHandler:  portSetHandler,
		settingsHandler: settingsHandler,
	}
}

func (s *MainScreen) Show() {
	//err := forms.RunSelectForm(func() (*forms.SelectFormSpec, error) {
	//	return &forms.SelectFormSpec{
	//		QuestionFn:     "K8s forwarder application",
	//		ErrorText: s.errorHandler.GetErrorText(),
	//		Items: []forms.OptionSpec{
	//			{
	//				Data: forms.OptionData{
	//					ID:          "forward",
	//					Name:        "Forward",
	//					Description: "Manage and enable port forwarding",
	//				},
	//				Func: func(data forms.OptionData) bool {
	//					if err := s.portSetHandler.Show(); err != nil {
	//						s.errorHandler.Handle(err, "failed to show port sets")
	//					}
	//
	//					return false
	//				},
	//			},
	//			{
	//				Data: forms.OptionData{
	//					ID:          "settings",
	//					Name:        "Settings",
	//					Description: "Update application settings",
	//				},
	//				Func: func(data forms.OptionData) bool {
	//					if err := s.settingsHandler.Show(); err != nil {
	//						s.errorHandler.Handle(err, "failed to show settings")
	//					}
	//
	//					return false
	//				},
	//			},
	//			{
	//				Data: forms.OptionData{
	//					ID:          "exit",
	//					Name:        promptui.Styler(promptui.FGItalic)("Exit"),
	//					Description: "Exit the application",
	//				},
	//				Func: func(data forms.OptionData) bool {
	//					return true
	//				},
	//			},
	//		},
	//	}, nil
	//})
	//if err != nil {
	//	s.errorHandler.Handle(err)
	//}

	var opts []selectpkg.Option
	for i := 0; i < 15; i++ {
		opts = append(opts, selectpkg.Option{
			Text: fmt.Sprintf("Forward %d", i),
			Desc: "Manage and enable port forwarding",
		})
	}

	cli.ClearScreen()
	selectModel := selectpkg.NewSelectModel(opts).
		SetHeaderFn(func() string { return "K8s forwarder application" }).
		SetQuestionFn(func() string { return "Choose an option" })

	p := tea.NewProgram(selectModel)
	model, err := p.Run()
	if err != nil {
		s.errorHandler.Handle(err)
	}

	cli.ClearScreen()
	fmt.Println(selectpkg.GetResultFromModel(model))

	time.Sleep(time.Second)

	return
}

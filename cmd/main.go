package main

import (
	"github.com/chzyer/readline"
	"github.com/yanodincov/k8s-forwarder/internal/config"
	"github.com/yanodincov/k8s-forwarder/internal/handler"
	forwardHandler "github.com/yanodincov/k8s-forwarder/internal/handler/forward"
	"github.com/yanodincov/k8s-forwarder/internal/handler/forward/services"
	settingsHandler "github.com/yanodincov/k8s-forwarder/internal/handler/settings"
	"github.com/yanodincov/k8s-forwarder/internal/handler/settings/files"
	"github.com/yanodincov/k8s-forwarder/internal/handler/settings/namespace"
	"github.com/yanodincov/k8s-forwarder/internal/repository/portset"
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/internal/service/forwarder"
	"github.com/yanodincov/k8s-forwarder/internal/service/k8s"
	"github.com/yanodincov/k8s-forwarder/pkg/cli"
	"github.com/yanodincov/k8s-forwarder/pkg/forms"
	"github.com/yanodincov/k8s-forwarder/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

func main() {
	readline.Stdout = cli.GetNoBellStdout()

	app := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return logger.NewFxLogger()
		}),

		// Config
		fx.Provide(config.NewConfig),

		// Tools
		fx.Provide(
			forms.NewErrorHandler,
		),
		fx.Invoke(forms.FixSurveyStyle),

		// Repositories
		fx.Provide(
			portset.NewStorage,
			portset.NewRepository,
			settings.NewStorage,
			settings.NewRepository,
		),

		// Services
		fx.Provide(
			k8s.NewService,
			forwarder.NewService,
		),

		// Screens
		fx.Provide(
			handler.NewMainScreen,

			forwardHandler.NewListScreen,
			forwardHandler.NewActionsScreen,
			forwardHandler.NewCreateScreen,
			forwardHandler.NewDeleteScreen,

			settingsHandler.NewListScreen,

			namespace.NewListScreen,
			namespace.NewActionListScreen,
			namespace.NewDeleteScreen,
			namespace.NewAddScreen,

			services.NewListScreen,
			services.NewAddScreen,
			services.NewActionsScreen,
			services.NewDeleteScreen,

			files.NewListScreen,
			files.NewAddScreen,
			files.NewActionListScreen,
			files.NewDeleteScreen,
		),

		// Run
		fx.Invoke(func(mainScreen *handler.MainScreen, exit fx.Shutdowner) {
			mainScreen.Show()
			cli.ClearScreen()
			_ = exit.Shutdown()
		}),
	)

	app.Run()
}

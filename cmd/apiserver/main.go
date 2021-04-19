package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"authkey/cmd/apiserver/internal/config"
	"authkey/cmd/apiserver/internal/handler"
	"authkey/cmd/apiserver/internal/svc"
	"authkey/pkg/lib"
	pconfig "authkey/pkg/lib/config"
	"authkey/pkg/lib/log"

	// "authkey/pkg/lib/ws"
	"authkey/pkg/util"

	"github.com/spf13/cobra"
	"github.com/unknwon/i18n"
	"go.uber.org/zap"
)

var (
	configFile      string
	survivalTimeout = 300 * time.Second

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "apiserver",
		Short: "",
		Run:   rootRun,
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version information",
		RunE: func(cmd *cobra.Command, _ []string) error {
			v := util.GetBuildInfo()
			fmt.Printf("%#v\n", v.String())
			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "f", `etc/apiserver*.ini`, "the config files")

	rootCmd.AddCommand(versionCmd)
}

func main() {
	rootCmd.Execute()
}

func rootRun(cmd *cobra.Command, args []string) {
	c := &config.Config{}
	pconfig.MustLoadINI(configFile, c)
	config.SetConfig(c)

	log.InitZap(c.Logger)
	defer log.Glog.Sync()

	log.Glog.Debug("config detail", zap.String("content", config.GlobalConfig.String()))
	log.Glog.Info("config loaded", zap.String("files", configFile))

	if config.GlobalConfig.Trace.Enable {
		shutdownFn, err := lib.InitOTEL(config.GlobalConfig.Trace.Server, config.GlobalConfig.Common.Name, log.Glog, log.GSlog)
		if err != nil {
			log.Glog.Fatal("init otel", zap.Error(err))
		}
		log.Glog.Info("init otel done")

		defer shutdownFn()
	}

	// load i18n
	if err := i18n.SetMessage("en", "etc/locale_en.ini"); err != nil {
		log.Glog.Fatal("load lang=en", zap.Error(err))
	}
	if err := i18n.SetMessage("zh", "etc/locale_zh.ini"); err != nil {
		log.Glog.Fatal("load lang=zh", zap.Error(err))
	}

	if err := svc.Init(); err != nil {
		log.Glog.Fatal("svc init", zap.Error(err))
	}

	// go ws.Manager.Run()
	go initSignal()

	h := handler.InitHandler2()
	if err := h.ListenAndServe(fmt.Sprintf("%s:%d", config.GlobalConfig.Common.Host, config.GlobalConfig.Common.Port)); err != nil {
		log.Glog.Fatal("start apiserver", zap.Error(err))
	}
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM)

	for {
		sig := <-signals
		log.Glog.Info("get signal", zap.String("signal", sig.String()))
		switch sig {
		case syscall.SIGHUP:
			// reload config
			c := &config.Config{}
			if err := pconfig.LoadINI(configFile, c); err == nil {
				config.SetConfig(c)
				log.Glog.Info("config reloaded", zap.String("files", configFile))
			} else {
				log.Glog.Error("config reloaded failed", zap.Error(err))
			}
		default:
			// The program exits normally or timeout forcibly exits.
			time.AfterFunc(time.Duration(survivalTimeout), func() {
				log.Glog.Error("exit force", zap.Duration("timeout", survivalTimeout))
				os.Exit(1)
			})

			log.Glog.Info("exit")
			os.Exit(0)
			return
		}
	}
}

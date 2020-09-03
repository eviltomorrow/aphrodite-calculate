package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"go.uber.org/zap"

	"github.com/eviltomorrow/aphrodite-base/zlog"
	"github.com/eviltomorrow/aphrodite-calculate/config"
)

const (
	nmConfig = "config"
)

var (
	path = flag.String(nmConfig, "config.toml", "配置文件路径")
)

var cfg = config.DefaultGlobalConfig
var cpf []func()

func main() {
	defer func() {
		if err := recover(); err != nil {
			zlog.Error("Panic: Unknown reason", zap.Error(fmt.Errorf("%v", err)))
			debug.PrintStack()
		}
		zlog.Sync()
	}()

	flag.Parse()

	if err := cfg.Load(*path, overrideFlags); err != nil {
		log.Printf("Warning: Load config file failure, use default config, nest error: %v\r\n", err)
	}

	setupLogConfig()
	setupGlobalVars()
	printInfo()
	registerCleanupFunc()
	blockingUntilTermination()

}

func setupLogConfig() {
	global, prop, err := zlog.InitLogger(&zlog.Config{
		Level:            cfg.Log.Level,
		Format:           cfg.Log.Format,
		DisableTimestamp: cfg.Log.DisableTimestamp,
		File: zlog.FileLogConfig{
			Filename: cfg.Log.FileName,
			MaxSize:  cfg.Log.MaxSize,
		},
		DisableStacktrace: true,
	})
	if err != nil {
		log.Printf("Fatal: Setup log config failure, nest error: %v", err)
		os.Exit(1)
	}
	zlog.ReplaceGlobals(global, prop)
}

func setupGlobalVars() {

}

func printInfo() {
	zlog.Info("Config information", zap.String("data", cfg.String()))
}

func overrideFlags(cfg *config.Config) {

}

func registerCleanupFunc() {

}

func blockingUntilTermination() {
	var ch = make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	switch <-ch {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	case syscall.SIGUSR1:
	case syscall.SIGUSR2:
	default:
	}
	for _, f := range cpf {
		f()
	}
	zlog.Info("Termination main programming, cleanup function is executed complete")
}

package main

import (
	"airpush/server"
	"bufio"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const APP_NAME = "rtb"

// load config
func initConfig() (v *viper.Viper, err error) {

	var path string

	flag.StringVar(&path, "config", "", "config path")
	flag.Parse()

	v = viper.New()
	v.SetEnvPrefix(APP_NAME)
	v.AutomaticEnv()
	v.SetConfigName("config")

	if path != "" {
		v.AddConfigPath(path)
	} else {
		v.AddConfigPath(fmt.Sprintf("$HOME/.%s", APP_NAME))
		v.AddConfigPath(".")
	}

	err = v.ReadInConfig()

	return
}

// init app settings
func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	rand.Seed(time.Now().UnixNano())
}

// loop app
func loop(exit func(os.Signal)) {
	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	switch sig := <-signalCh; sig {
	case os.Interrupt:
		exit(sig)
	case syscall.SIGTERM:
		exit(sig)
	case syscall.SIGQUIT:
		exit(sig)
	}
}

// simple rtb app
func main()  {


	// init logger
	logger := logrus.New()

	// init config
	config, err := initConfig()
	if err != nil {
		logger.Fatal("load config fail, err: %s", err)
	} else {

		// setting how max app used cpu core
		if c := config.GetInt("app.setting.cpu_core"); c != 0 {
			runtime.GOMAXPROCS(config.GetInt("app.setting.cpu_core"))
		} else {
			runtime.GOMAXPROCS(runtime.NumCPU())
		}

		// setting logger level
		// logger have many hooks graylog/telegram/etc... I think is nice
		if c := config.GetString("app.setting.mode"); c != "debug" {
			logger.SetOutput(bufio.NewWriterSize(os.Stdout, 1024*16)) // 16kb pre buffer output
			logger.SetLevel(logrus.ErrorLevel)
		} else {
			logger.SetOutput(os.Stdout)
			logger.SetLevel(logrus.TraceLevel)
		}
	}

	// init server
	s, err := server.New(

		server.SetConcurrency(config.GetInt("app.server.Concurrency")),
		server.SetDisableKeepalive(config.GetBool("app.server.DisableKeepalive")),

		server.SetReadBufferSize(config.GetInt("app.server.ReadBufferSize")),
		server.SetWriteBufferSize(config.GetInt("app.server.WriteBufferSize")),

		server.SetWriteTimeout(config.GetInt("app.server.WriteTimeout")),
		server.SetReadTimeout(config.GetInt("app.server.ReadTimeout")),

		server.SetServerName("simple rtb"),
		server.SetServerAddr(config.GetString("app.server.ServerAddr")),
		server.SetLogger(logger),
	)
	if err != nil {
		logger.Fatalf("init server fail: %s", err)
	}

	// start server
	go func() {
		err = s.Start()
		if err != nil {
			logger.Fatalf("start server fail: %s", err)
		}
	}()

	logger.Infof("app loaded with conf: %s", config.ConfigFileUsed())

	// loop app
	loop(func(i os.Signal) {
		logger.Info("graceful shutdown...")
		_ = s.Close()
	})
}


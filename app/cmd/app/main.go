package main

import (
	"app/internal/config"
	delivery "app/internal/delivery/http"

	"fmt"
	"os"
	"os/signal"
	"syscall"

	lg "github.com/dmitryavdonin/gtools/logger"
)

var Version = "2.0.0"

func main() {
	cfg, err := config.InitConfig("")
	if err != nil {
		panic(fmt.Sprintf("error initializing config %s", err))
	}

	//setup logger
	logger, err := lg.New(cfg.Log.Level, cfg.App.ServiceName)
	if err != nil {
		panic(fmt.Sprintf("error initializing logger %s", err))
	}

	delivery, err := delivery.New(Version, cfg.App.Port, logger, delivery.Options{})
	if err != nil {
		logger.Fatal("delivery initialization error: %s", err.Error())
	}

	//logger.Info("main(): Auth app version = " + Version)
	fmt.Println("main(): app version = " + Version)

	err = delivery.Run()
	if err != nil {
		logger.Fatal("start delivery error: %s", err.Error())
	}

	//closes connections on app kill
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}

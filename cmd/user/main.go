package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/grozaqueen/julse/internal/apps/user"
	"github.com/grozaqueen/julse/internal/configs"
	"github.com/grozaqueen/julse/internal/configs/logger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

const (
	userService = "user_go"
	configFile  = ".env"
	kafka       = "kafka"
)

// todo вынос в apps
func main() {
	err := godotenv.Load(configFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	viper, err := configs.SetupViper()
	if err != nil {
		log.Fatalf("Error setting up viper %v", err)
	}

	serviceConf := viper.GetStringMap(userService)
	kafkaConf := viper.GetStringMap(kafka)

	slogLog := logger.InitLogger()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor())

	app, err := user.NewUsersApp(slogLog, server, serviceConf, kafkaConf)

	router := mux.NewRouter()
	router.PathPrefix("/metrics").Handler(promhttp.Handler())
	serverProm := http.Server{Handler: router, Addr: fmt.Sprintf(":%d", 8081), ReadHeaderTimeout: 10 * time.Second}

	go func() {
		if err = serverProm.ListenAndServe(); err != nil {
			log.Println("fail auth.ListenAndServe")
		}
	}()

	err = app.Run()
	if err != nil {
		log.Fatalf("err %v", err)
	}
}

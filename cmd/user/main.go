package main

import (
	"github.com/grozaqueen/julse/internal/apps/user"
	"github.com/grozaqueen/julse/internal/configs"
	"github.com/grozaqueen/julse/internal/configs/logger"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
)

const (
	userService = "user_go"
	configFile  = ".env"
)

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

	slogLog := logger.InitLogger()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor())

	app, err := user.NewUsersApp(slogLog, server, serviceConf)

	err = app.Run()
	if err != nil {
		log.Fatalf("err %v", err)
	}
}

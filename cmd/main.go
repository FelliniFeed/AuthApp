package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/FelliniFeed/AuthApp.git/pkg"
	"github.com/FelliniFeed/AuthApp.git/pkg/handler"
	"github.com/FelliniFeed/AuthApp.git/pkg/repository"
	"github.com/FelliniFeed/AuthApp.git/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := gotenv.Load("../.env"); err != nil {
		logrus.Fatalf("error initializing envs: %s", err.Error())
		return
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(pkg.Server)

	go func() {
		if err := srv.Run(viper.GetString("server.port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Server started on", viper.GetString("server.port"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Print("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("failed to shutdown server: %s", err)
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("failed to close database: %s", err)
	}

}

func initConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
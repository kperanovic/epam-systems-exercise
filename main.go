package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	middleware "github.com/kperanovic/epam-systems/api/v1/auth"
	"github.com/kperanovic/epam-systems/api/v1/handlers"
	"github.com/kperanovic/epam-systems/internal/kafka"
	"github.com/kperanovic/epam-systems/internal/logger"
	"github.com/kperanovic/epam-systems/internal/storage"
	"github.com/kperanovic/epam-systems/internal/token"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	if err := loadParams(); err != nil {
		panic(err)
	}

	log := logger.NewProduction()

	log.Info("starting service")

	store := storage.NewMySQLStorage()
	if err := store.Connect(); err != nil {
		log.Fatal("error establishing connection", zap.Error(err))
	}

	producer, err := kafka.NewKafkaProducer(
		viper.GetStringSlice("KAFKA_ADDR"),
		viper.GetString("KAFKA_TOPIC"),
		log,
	)
	if err != nil {
		log.Fatal("error starting kafka producer", zap.Error(err))
	}

	h := handlers.NewRESTHandlers(log, store, producer)

	r := gin.Default()

	t, err := token.NewJWTToken(viper.GetString("AUTH_SECRET"))
	if err != nil {
		log.Fatal("error creating jwt token instance", zap.Error(err))
	}

	group := r.Group("v1/company").Use(middleware.AuthMiddleware(t))
	group.POST("/", h.HandleCreateCompany)
	group.PATCH("/:id", h.HandlePatchCompany)
	group.DELETE("/:id", h.HandleDeleteCompany)

	r.GET("/v1/company/:id", h.HandleGetCompany)

	if err := r.Run(); err != nil {
		log.Fatal("error starting http server", zap.Error(err))
	}
}

func loadParams() error {
	mandatory := []string{
		"AUTH_SECRET",
		"KAFKA_ADDR",
		"DB_USER",
		"DB_PWD",
	}

	viper.SetDefault("DB_HOST", "127.0.0.1:3306")
	viper.SetDefault("DB_NAME", "epam")
	viper.SetDefault("KAFKA_VERSION", "1.0.0")

	viper.AutomaticEnv()

	for _, param := range mandatory {
		if !viper.IsSet(param) {
			return fmt.Errorf("mandatory parameters not set (%s)", param)
		}
	}

	return nil
}

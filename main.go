package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"

	"github.com/pawverse/pawcare-core/pkg/common"
	"github.com/pawverse/pawcare-core/pkg/db/mongodb"
	"github.com/pawverse/pawcare-medical/internal/config"
	"github.com/pawverse/pawcare-medical/internal/pet/endpoint"
	"github.com/pawverse/pawcare-medical/internal/pet/service"
	recordendpoint "github.com/pawverse/pawcare-medical/internal/record/endpoint"
	recordservice "github.com/pawverse/pawcare-medical/internal/record/service"
	"github.com/pawverse/pawcare-medical/internal/repository/mongo"
	"github.com/pawverse/pawcare-medical/internal/transport"

	"github.com/joho/godotenv"
)

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func _main() error {
	godotenv.Load()

	viper := config.InitConfig()
	ctx := context.Background()

	http.DefaultTransport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = viper.GetBool(common.InsecureSkipVerifyKey)

	db, teardown, err := mongodb.ConnectDB(ctx, viper.GetString(common.DBConnectionStringKey), "medical")
	defer teardown(ctx)
	if err != nil {
		return err
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	petRepository := mongo.NewPetRepository(db)
	petService := service.NewPetService(petRepository, logger.With(zap.String("service", "pet")))
	petEndpoints, err := endpoint.NewSet(viper, petService)
	if err != nil {
		return err
	}

	router, err := transport.NewRouter(viper, petEndpoints, logger.With(zap.String("transport", "kafka")))
	if err != nil {
		return err
	}

	repository := mongo.NewRecordRepository(db, logger.With(zap.String("repository", "record")))
	recordService := recordservice.NewRecordService(repository, petService, logger.With(zap.String("service", "record")))
	recordEndpoints, err := recordendpoint.NewSet(viper, recordService)
	if err != nil {
		return err
	}

	httpHandler := transport.MakeHTTPServer(recordEndpoints, logger.With(zap.String("transport", "http")))
	httpHandler = accessControl(httpHandler)

	var g group.Group

	{
		httpAddr := fmt.Sprintf(":%s", viper.GetString(common.HTTPPortKey))
		logger := logger.With(zap.String("transport", "http"))
		httpListenAddr, err := net.Listen("tcp", httpAddr)
		if err != nil {
			return err
		}

		g.Add(func() error {
			logger.Info("Starting server", zap.String("addr", httpAddr))
			return http.Serve(httpListenAddr, httpHandler)
		}, func(err error) {
			logger.Info("Closing server", zap.NamedError("reason", err))
			httpListenAddr.Close()
		})
	}

	g.Add(func() error {
		logger.Info("Starting message router")
		return router.Run(ctx)
	}, func(err error) {
		logger.Info("Closing message router", zap.NamedError("reason", err))
		router.Close()
	})

	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		return fmt.Errorf("%s", <-c)
	}, func(err error) {
		logger.Info("Shutdown signal received", zap.NamedError("signal", err))
	})

	return g.Run()
}

func main() {
	if err := _main(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

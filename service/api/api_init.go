package api

import (
	"example_project/internal/config"
	"example_project/internal/database"
	"example_project/internal/leveledlog"
	"example_project/internal/redisclient"
	"example_project/internal/server"
	"example_project/internal/version"
	"example_project/service/api/utils"
	"os"
	"sync"

	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator"
)

type CoreAppStruct struct {
	utils.Application
}

func Serve() {
	cfg, err := config.SetupConfig()
	if err != nil {
		panic(err)
	}

	logger := leveledlog.NewLogger(os.Stdout, leveledlog.GetLevel(cfg.LogLevel), true).WithID("main")

	logger.Info("master database url: %s", cfg.MasterDbHost)
	db, err := database.New(cfg.MasterDbUrl, false)
	if err != nil {
		logger.Fatal(err.Error())
	}

	defer db.Close()

	app := &CoreAppStruct{
		Application: utils.Application{
			Config:       *cfg,
			Db:           db,
			Logger:       logger,
			Validator:    validator.New(),
			Wg:           &sync.WaitGroup{},
			TokenManager: jwtauth.New("HS256", []byte(cfg.SecretKey), nil),
		},
	}

	// this wg was added as server is in running mode if its closed then we need to call app.wg.done()
	app.Wg.Add(1)
	consumerDone := make(chan bool)
	consumer := redisclient.Redis{Config: *cfg, Db: db, Logger: logger, Wg: app.Wg, DistributedRedisClient: nil}

	go consumer.ConsumerEvent(consumerDone, app.Wg)
	logger.Info("starting api server on %s (version: %s)",
		cfg.GetServerAddress(), version.GetVersion())
	err = server.Run(cfg.GetServerAddress(), app.routes())
	app.Wg.Done()
	if err != nil {
		logger.Error("error form the server--------->")
		logger.Fatal(err.Error())
	}
	logger.Info("server stopped")
	// go utils.HandleUnprocessedEvent(retryWorkerContext, app.Wg, time.Duration(app.Config.RetryCount))

	// retryWorkerContextCancel() // cancel retry worker context
	logger.Info("closed consumer routine, waiting on application waitgroup")
	logger.Info("exiting application, graceful shutdown completed")
	<-consumerDone

}

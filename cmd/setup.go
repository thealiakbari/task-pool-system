package cmd

import (
	"context"
	taskHttpAdaptor "github.com/thealiakbari/task-pool-system/internal/adapters/inbound/http/task"
	taskOutboundRepo "github.com/thealiakbari/task-pool-system/internal/adapters/outbound/db/pg"
	taskApp "github.com/thealiakbari/task-pool-system/internal/application/task"
	taskService "github.com/thealiakbari/task-pool-system/internal/domain/task"
	"github.com/thealiakbari/task-pool-system/internal/domain/task/pool"
	taskInterface "github.com/thealiakbari/task-pool-system/internal/ports/inbound/task"
	taskRepo "github.com/thealiakbari/task-pool-system/internal/ports/outbound/task"
	"github.com/thealiakbari/task-pool-system/pkg/common/config"
	"github.com/thealiakbari/task-pool-system/pkg/common/db"
	"github.com/thealiakbari/task-pool-system/pkg/common/i18next"
	"github.com/thealiakbari/task-pool-system/pkg/common/logger"
	"golang.org/x/text/language"
)

type RepositoryStorage struct {
	taskRepo taskRepo.TaskRepository
}

type ServiceStorage struct {
	taskSvc taskInterface.TaskService
}

type ApplicationStorage struct {
	taskApp taskApp.TaskHttpApp
}

type HttpAdaptorStorage struct {
	TaskAdaptor taskHttpAdaptor.Adaptor
}

type SetupConfig struct {
	Ctx                context.Context
	Conf               *config.AppConfig
	Logger             logger.Logger
	DB                 db.DBWrapper
	HttpAdaptorStorage HttpAdaptorStorage
}

func Setup() *SetupConfig {
	ctx := context.Background()
	conf := config.LoadConfig("./config/config.yml")

	log, err := logger.New(
		conf.Mode,
		conf.ServiceName,
		"todoapp",
	)
	if err != nil {
		panic(err)
	}

	err = i18next.NewLanguage(language.Make(conf.Language))
	if err != nil {
		panic(err)
	}

	logInfra := log.CloneAsInfra()
	err = db.Migrate(conf.DB.Postgres, logInfra)
	if err != nil {
		logInfra.Panicf("Migration failed: %s\n", err.Error())
	}
	logInfra.Info("Migrations successfully done.")

	gormDB, err := db.NewPostgresConn(ctx, conf.DB.Postgres)
	if err != nil {
		panic(err)
	}

	dbw := db.NewDBWrapper(gormDB)

	repos := NewRepositoryStorage(dbw)
	services := NewServiceStorage(log, repos)

	httpApps := NewHttpAppStorage(dbw, services)
	httpAdaptors := NewHttpAdaptorStorage(httpApps)

	return &SetupConfig{
		Ctx:                ctx,
		Conf:               conf,
		Logger:             log,
		DB:                 dbw,
		HttpAdaptorStorage: httpAdaptors,
	}
}

func NewHttpAppStorage(
	db db.DBWrapper,
	services ServiceStorage,
) ApplicationStorage {
	poolWorker := pool.New(context.Background(), 10, 10)
	poolWorker.Start(pool.WorkerDeps{
		TaskService: services.taskSvc,
	})
	return ApplicationStorage{
		taskApp: taskApp.NewTaskHttpApp(services.taskSvc, db, poolWorker),
	}
}

func NewRepositoryStorage(db db.DBWrapper) RepositoryStorage {
	return RepositoryStorage{
		taskRepo: taskOutboundRepo.NewTaskRepository(db),
	}
}

func NewServiceStorage(log logger.Logger, repos RepositoryStorage) ServiceStorage {
	taskSvc := taskService.NewTaskService(taskService.TaskConfig{Logger: log, TaskRepo: repos.taskRepo})

	return ServiceStorage{
		taskSvc: taskSvc,
	}
}

func NewHttpAdaptorStorage(
	httpApps ApplicationStorage,
) HttpAdaptorStorage {
	return HttpAdaptorStorage{
		TaskAdaptor: taskHttpAdaptor.Adaptor{TaskHttpApp: httpApps.taskApp},
	}
}

package server

import (
	"context"
	"github.com/ulysseskk/common/components"
	"github.com/ulysseskk/common/components/redis"
	"github.com/ulysseskk/common/components/sql"
	"github.com/ulysseskk/common/config"
	"github.com/ulysseskk/common/health"
	"github.com/ulysseskk/common/logger/log"
	"github.com/ulysseskk/common/model/errors"
	"github.com/ulysseskk/common/trace"
	"github.com/ulysseskk/common/util/goroutineUtil"
	"os"
	"os/signal"
	"syscall"
)

type IServer interface {
	Init() error
	Start() error
	Shutdown() error
}

type BaseServer[T IBaseConfig] struct {
	cancelCtx     context.Context
	cancelFunc    context.CancelFunc
	serviceConfig T
}

func (b *BaseServer[T]) Ctx() context.Context {
	return b.cancelCtx
}

func (b *BaseServer[T]) DoStart(startFunc func(ctx context.Context) error) error {
	errch := make(chan error, 1)
	go goroutineUtil.SafeGoroutineWithLog(func() {
		err := startFunc(b.cancelCtx)
		if err != nil {
			errch <- err
		}
	})
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGPIPE, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)

	exitCode := 0
Loop:
	for {
		select {
		case <-b.cancelCtx.Done():
			log.GlobalLogger().Infof("controller is exiting as %s", b.cancelCtx.Err())
			break Loop
		case sig := <-signals:
			log.GlobalLogger().Infof("received signal: %s", sig)
			switch sig {
			case syscall.SIGPIPE:
			case syscall.SIGINT:
				// SIGINT is used for the notification of config changed.
				// Set exit code = 1 to let process manager tool (like systemctl) restart current process.
				exitCode = 1
				fallthrough
			default:
				log.GlobalLogger().Infof("controller  is exiting as received signal: %s", sig)
				break Loop
			}
		case err := <-errch:
			if err != nil {
				log.GlobalLogger().Errorf("controller  is exiting as failed to start one of services: %s", err)
				break Loop
			}
		}
	}
	b.cancelFunc()
	if exitCode != 0 {
		os.Exit(exitCode)
	}
	return nil
}

func (b *BaseServer[T]) Init(conf T) error {
	b.cancelCtx, b.cancelFunc = context.WithCancel(context.Background())
	log.GlobalLogger().Debug("init base server")
	err := config.InitConfig[T](conf)
	if err != nil {
		return errors.NewError().WithError(err).WithCode(errors.CodeInitializeError).WithMessage("初始化配置失败")
	}
	health.AddDefaultRegister("/configs/static", func() (interface{}, error) {
		return conf, nil
	})
	if err := conf.Validate(); err != nil {
		return err
	}
	b.serviceConfig = conf
	log.GlobalLogger().Debug("init components")
	err = b.initComponents()
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseServer[T]) initComponents() error {
	// 初始化Logger
	if b.serviceConfig.GetLoggerConfig() != nil {
		log.GlobalLogger().Debug("init logger")
		err := log.InitGlobalLogger(b.serviceConfig.GetLoggerConfig())
		if err != nil {
			return errors.NewError().WithError(err).WithCode(errors.CodeInitializeError).WithMessage("初始化Logger失败")
		}
		log.GlobalLogger().Debug("init logger success")
	}
	// 初始化Tracer
	log.GlobalLogger().Debug("init tracer")
	err := trace.InitTracer(b.serviceConfig.AppConfig().Name)
	if err != nil {
		return errors.NewError().WithError(err).WithCode(errors.CodeInitializeError).WithMessage("初始化Tracer失败")
	}
	// 如果存在初始化DB
	if b.serviceConfig.GetSqlConfig() != nil {
		log.GlobalLogger().Debug("init db")
		_, err := sql.InitDefault(*b.serviceConfig.GetSqlConfig(), sql.WithRestErrorStackCallback())
		if err != nil {
			return errors.NewError().WithError(err).WithCode(errors.CodeInitializeError).WithMessage("初始化DB失败")
		}
	}
	if b.serviceConfig.ComponentsConfig() != nil {
		log.GlobalLogger().Debug("init components")
		err := components.InitComponents(b.serviceConfig.ComponentsConfig())
		if err != nil {
			return err
		}
	}
	// 初始化多db
	if b.serviceConfig.MultiDatabaseConfig() != nil {
		log.GlobalLogger().Debug("init multi db")
		err := sql.InitMulti(b.serviceConfig.MultiDatabaseConfig(), sql.WithRestErrorStackCallback())
		if err != nil {
			return errors.NewError().WithError(err).WithCode(errors.CodeInitializeError).WithMessage("初始化多DB失败")
		}
		log.GlobalLogger().Debug("init multi db success")
	}
	if b.serviceConfig.RedisConfig() != nil {
		log.GlobalLogger().Debug("init redis")
		err = redis.InitRedis(b.serviceConfig.RedisConfig())
		if err != nil {
			return errors.NewError().WithError(err).WithCode(errors.CodeInitializeError).WithMessage("初始化Redis失败")
		}
	}
	// 初始化Health
	if b.serviceConfig.HealthConfig() != nil {
		log.GlobalLogger().Debug("init health")
		health.InitHealthServer(b.serviceConfig.HealthConfig())
	}
	return nil
}

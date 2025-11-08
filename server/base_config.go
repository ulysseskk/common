package server

import (
	"fmt"
	"time"

	"github.com/ulysseskk/common/components"
	"github.com/ulysseskk/common/components/redis"
	"github.com/ulysseskk/common/components/sql"
	"github.com/ulysseskk/common/health"
	loggerConfig "github.com/ulysseskk/common/logger/conf"
)

type IBaseConfig interface {
	GetSqlConfig() *sql.DatabaseConfig
	GetLoggerConfig() *loggerConfig.LogConfig
	HttpServerConfig() *GinConfig
	AppConfig() *AppConfig
	HealthConfig() *health.Config
	RedisConfig() *redis.Config
	MultiDatabaseConfig() sql.MultiDatabaseConfig
	ComponentsConfig() *components.Config
	Validate() error
}

// Common 组合这个struct会让所有的配置都在common下，不要随便改名，不然所有配置的根path都要改
type Common struct {
	App           *AppConfig              `json:"app" yaml:"app"`
	Database      *sql.DatabaseConfig     `json:"database" yaml:"database"`
	Logger        *loggerConfig.LogConfig `json:"logger" yaml:"logger"`
	Http          *GinConfig              `json:"http" yaml:"http"`
	Health        *health.Config          `json:"health" yaml:"health"`
	Redis         *redis.Config           `json:"redis" yaml:"redis"`
	MultiDatabase sql.MultiDatabaseConfig `json:"multi_database" yaml:"multi_database"`
	Components    *components.Config      `json:"components" yaml:"components"`
}

func (b Common) GetSqlConfig() *sql.DatabaseConfig {
	return b.Database
}

func (b Common) GetLoggerConfig() *loggerConfig.LogConfig {
	return b.Logger
}

func (b Common) HttpServerConfig() *GinConfig {
	return b.Http
}

func (b Common) AppConfig() *AppConfig {
	return b.App
}

func (b Common) HealthConfig() *health.Config {
	return b.Health
}

func (b Common) RedisConfig() *redis.Config {
	return b.Redis
}

func (b Common) MultiDatabaseConfig() sql.MultiDatabaseConfig {
	return b.MultiDatabase
}

func (b Common) ComponentsConfig() *components.Config {
	return b.Components
}

func (b Common) Validate() error {
	if b.App == nil {
		return fmt.Errorf("未发现App配置，请补充！")
	}
	return nil
}

type AppConfig struct {
	Name       string           `json:"name"`      // 服务名
	Env        string           `json:"env"`       // 环境
	Namespace  string           `json:"namespace"` // 命名空间
	Controller ControllerConfig `json:"controller"`
}

type ControllerConfig struct {
	EnableLeaderElection bool `json:"enable_leader_election" yaml:"enable_leader_election"`
}

type GinConfig struct {
	Listen                  *ListenConfig `json:"listen"`
	GracefulShutdownTimeout time.Duration `json:"graceful_shutdown_timeout"`
	Mode                    string        `json:"mode"`
}

type ListenConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

package config

import (
	"gitlab.ulyssesk.top/common/common/logger/log"
	"fmt"
	"os"
)

var (
	defaultLoaderConfig *LoaderConfig
	ConfigBasePath      string
)

type LoaderConfig struct {
	LoaderType string                   `json:"loader_type"`
	Local      *LocalConfigLoaderConfig `json:"local"`
	Gzip       *GzipConfigLoaderConfig  `json:"gzip"`
}

const (
	LoaderTypeEmpty     = "empty"
	LoaderTypeApollo    = "apollo"
	LoaderTypeLocal     = "local"
	LoaderTypeConfigMap = "config_map"
	LoaderTypeGzip      = "gzip"
)

func init() {
	if os.Getenv("GZIP_CONFIG_PATH") != "" {
		gzipPath := os.Getenv("GZIP_CONFIG_PATH")
		targetPath := os.Getenv("GZIP_CONFIG_TARGET_PATH")
		if targetPath == "" {
			panic("GZIP_CONFIG_TARGET_PATH is required")
		}
		defaultLoaderConfig = &LoaderConfig{
			LoaderType: LoaderTypeGzip,
			Gzip: &GzipConfigLoaderConfig{
				GzipPath:   gzipPath,
				TargetPath: targetPath,
				FileType:   "yaml",
			},
		}
	} else {
		defaultLoaderConfig = &LoaderConfig{
			LoaderType: LoaderTypeLocal,
			Local: &LocalConfigLoaderConfig{
				ConfigPath: "config.yaml",
				Absolute:   false,
				FileType:   "yaml",
			},
		}
		if os.Getenv("CONFIG_LOADER_TYPE") != "" {
			defaultLoaderConfig.LoaderType = os.Getenv("CONFIG_LOADER_TYPE")
		}
		if os.Getenv("CONFIG_LOADER_LOCAL_CONFIG_PATH") != "" {
			defaultLoaderConfig.Local.ConfigPath = os.Getenv("CONFIG_LOADER_LOCAL_CONFIG_PATH")
			defaultLoaderConfig.Local.Absolute = isAbsolute(defaultLoaderConfig.Local.ConfigPath)
		}
		if os.Getenv("CONFIG_LOADER_LOCAL_CONFIG_FILE_TYPE") != "" {
			defaultLoaderConfig.Local.FileType = os.Getenv("CONFIG_LOADER_LOCAL_CONFIG_FILE_TYPE")
		}
	}

}

func InitConfig[T any](configStructPointer T) error {
	return InitConfigWithLoaderConfig(*defaultLoaderConfig, configStructPointer)
}

func InitConfigWithLoaderConfig[T any](conf LoaderConfig, configStructPointer T) error {
	log.GlobalLogger().Debugf("init config with loader config: %+v", conf.LoaderType)
	var loader Loader
	switch conf.LoaderType {
	case LoaderTypeApollo:
		return nil
	case LoaderTypeConfigMap:
		return fmt.Errorf("not implemented yet")
	case LoaderTypeLocal:
		loader = &LocalConfigLoader{*conf.Local}
		err := loader.Init(configStructPointer)
		if err != nil {
			return err
		}
	case LoaderTypeGzip:
		loader = &GzipConfigLoader{conf: *conf.Gzip}
		err := loader.Init(configStructPointer)
		if err != nil {
			return err
		}
	case LoaderTypeEmpty:
		return nil
	default:
		return fmt.Errorf("unsuportted loader type %s", conf.LoaderType)
	}
	log.GlobalLogger().Debug("init config success")
	return nil
}

func isAbsolute(path string) bool {
	return path[0] == '/' || path[0] == '\\'
}

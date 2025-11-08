package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"gitlab.ulyssesk.top/common/common/logger/log"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	fileNotExistError = fmt.Errorf("config file not exist")
	viperInstance     = viper.New()
)

type LocalConfigLoaderConfig struct {
	ConfigPath string `json:"config_path"`
	Absolute   bool   `json:"absolute"`
	FileType   string `json:"file_type"`
}

type LocalConfigLoader struct {
	conf LocalConfigLoaderConfig
}

func (l *LocalConfigLoader) Init(configStructPointer any) error {
	log.GlobalLogger().Debug("init local config loader")
	// 先拼文件路径
	configFilePath := ""
	if !l.conf.Absolute {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		configFilePath = filepath.Join(wd, l.conf.ConfigPath)
	} else {
		configFilePath = l.conf.ConfigPath
	}
	log.GlobalLogger().Debugf("config file path: %s.Start stat configfile", configFilePath)
	if _, err := fs.Stat(os.DirFS("/"), configFilePath); err != nil && errors.Is(err, fs.ErrNotExist) {
		return fileNotExistError
	}
	log.GlobalLogger().Debug("config file exist.Read config file")
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	log.GlobalLogger().Debug("read config file success")
	viperInstance.SetConfigType(l.conf.FileType)
	err = viperInstance.ReadConfig(bytes.NewBuffer(fileContent))
	if err != nil {
		return err
	}
	err = viperInstance.Unmarshal(configStructPointer, func(config *mapstructure.DecoderConfig) {
		config.TagName = l.conf.FileType
	})
	if err != nil {
		return err
	}
	log.GlobalLogger().Debug("unmarshal config success")
	return nil
}

func (l *LocalConfigLoader) SetCallback(callback ReloadConfigCallback) {
	//TODO implement me
	panic("implement me")
}

func GlobalViperInstance() *viper.Viper {
	return viperInstance
}

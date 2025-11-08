package config

import (
	"bytes"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"gitlab.ulyssesk.top/common/common/util/compress"
	"os"
	"path/filepath"
)

type GzipConfigLoaderConfig struct {
	GzipPath   string `json:"gzip_path"`
	TargetPath string `json:"target_path"`
	FileType   string `json:"file_type"`
}

type GzipConfigLoader struct {
	conf GzipConfigLoaderConfig
}

func (g GzipConfigLoader) Init(configStructPointer any) error {
	compressPath := g.conf.TargetPath
	err := os.MkdirAll(compressPath, os.ModePerm)
	if err != nil {
		return err
	}
	configFilePath := filepath.Join(compressPath, "main.yaml")

	// 先看是否有Gzip文件
	err = compress.DecompressTarGz(g.conf.GzipPath, compressPath)
	if err != nil {
		return err
	}
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	v := viper.New()
	v.SetConfigType(g.conf.FileType)
	err = v.ReadConfig(bytes.NewBuffer(fileContent))
	if err != nil {
		return err
	}
	err = v.Unmarshal(configStructPointer, func(config *mapstructure.DecoderConfig) {
		config.TagName = g.conf.FileType
	})
	if err != nil {
		return err
	}
	ConfigBasePath = compressPath
	return nil
}

func (g GzipConfigLoader) SetCallback(callback ReloadConfigCallback) {
}

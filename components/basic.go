package components

import (
	"github.com/ulysseskk/common/components/minio"
	"github.com/ulysseskk/common/components/opensearch"
	"github.com/ulysseskk/common/logger/log"
	"github.com/ulysseskk/common/model/errors"
)

type Config struct {
	Opensearch *opensearch.Config `json:"opensearch" yaml:"opensearch"`
	Minio      *minio.Config      `json:"minio" yaml:"minio"`
}

func InitComponents(cfg *Config) error {
	if cfg.Opensearch != nil {
		log.GlobalLogger().Debug("Init opensearch client")
		err := opensearch.Init(cfg.Opensearch)
		if err != nil {
			log.GlobalLogger().Errorf("Fail to init opensearch client: %v", err)
			return errors.NewError().WithError(err).WithCode(errors.CodeInitializeError).WithMessage("fail to init opensearch client")
		}
		log.GlobalLogger().Debug("Init opensearch client success")
	}
	if cfg.Minio != nil {
		log.GlobalLogger().Debug("Init minio client")
		minio.Init(cfg.Minio)
		log.GlobalLogger().Debug("Init minio client success")
	}
	return nil
}

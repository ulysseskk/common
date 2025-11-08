package opensearch

import (
	"crypto/tls"
	"fmt"
	"github.com/opensearch-project/opensearch-go"
	"github.com/ulysseskk/common/model/errors"
	"net/http"
)

type Clients struct {
	cfg              *Config
	opensearchClient *opensearch.Client
}

type Config struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	UserName string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func NewClients(cfg *Config) (*Clients, error) {
	var err error
	actualCfg := &Config{
		Host:     cfg.Host,
		Port:     cfg.Port,
		UserName: cfg.UserName,
		Password: cfg.Password,
	}
	cli := &Clients{
		cfg: actualCfg,
	}
	err = cli.init()
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (c *Clients) init() error {
	var err error
	c.opensearchClient, err = opensearch.NewClient(opensearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://%s:%d", c.cfg.Host, c.cfg.Port),
		},
		Username: c.cfg.UserName,
		Password: c.cfg.Password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	if err != nil {
		return errors.NewError().WithError(err).WithCode(errors.CodeInitializeError).WithMessage("fail to create opensearch client")
	}
	return nil
}

func (c *Clients) GetOpenSearchClient() *opensearch.Client {
	return c.opensearchClient
}

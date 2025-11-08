package opensearch

import (
	"github.com/opensearch-project/opensearch-go"
)

var (
	c *Clients
)

func Init(cfg *Config) error {
	var err error
	c, err = NewClients(cfg)
	if err != nil {
		return err
	}
	return nil
}

func GetOpenSearchClient() *opensearch.Client {
	return c.GetOpenSearchClient()
}

func GetSingleton() *Clients {
	return c
}

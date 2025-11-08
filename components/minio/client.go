package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Clusters map[string]ClusterConfig `json:"clusters" yaml:"clusters"`
}

type ClusterConfig struct {
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	AccessKey string `json:"accessKey" yaml:"accessKey"`
	SecretKey string `json:"secretKey" yaml:"secretKey"`
	Token     string `json:"token" yaml:"token"`
	Secure    bool   `json:"secure" yaml:"secure"`
	Region    string `json:"region" yaml:"region"`
}

var (
	config   *Config
	clusters map[string]*minio.Client
)

func Init(c *Config) {
	config = c
	clusters = make(map[string]*minio.Client)
	for name, cluster := range config.Clusters {
		clusters[name], _ = minio.New(cluster.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(cluster.AccessKey, cluster.SecretKey, cluster.Token),
			Secure: cluster.Secure,
			Region: cluster.Region,
		})
	}
}

func GetCluster(name string) *minio.Client {
	return clusters[name]
}

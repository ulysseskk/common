package redis

type Config struct {
	Mode       string   `json:"mode" yaml:"mode"`
	Cluster    bool     `json:"cluster" yaml:"cluster"`
	Addrs      []string `json:"addrs" yaml:"addrs"`
	User       string   `json:"user" yaml:"user"`
	Password   string   `json:"password" yaml:"password"`
	MasterName string   `json:"master_name" yaml:"master_name"`
}

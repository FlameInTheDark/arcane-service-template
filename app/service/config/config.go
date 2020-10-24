package config

import (
	"fmt"
	"strings"
	"sync"
)

const (
	EtcdDatabase    = "/services/arcane/config/database"
	EtcdNats        = "/conf/nats/endpoints"
	EtcdEnvironment = "/services/arcane/config/environment"
	EtcdMetrics     = "/services/arcane/config/metrics"
	EtcdDiscord     = "/service/arcane/config/discord"
)

type Service struct {
	sync.Mutex
	Environment EnvironmentConfig
	Nats        NatsConfig
	Database    DatabaseConfig
	Metrics     MetricsConfig
	Discord     DiscordConfig
}

type EnvironmentConfig struct {
	Discord  string `json:"discord"`
	Database string `json:"database"`
}

type DiscordConfig struct {
	Token         string `json:"token"`
	OAuthClientID string `json:"oauth_client_id"`
	OAuthSecret   string `json:"oauth_secret"`
	OAuthScope    string `json:"oauth_scope"`
	OAuthRedirect string `json:"oauth_redirect"`
	InviteURL     string `json:"invite"`
}

type NatsConfig []NatsEndpoint

type NatsEndpoint struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type DatabaseConfig []DatabaseEndpoint

type DatabaseEndpoint struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type MetricsConfig struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
	Org      string `json:"org"`
	Bucket   string `json:"bucket"`
}

func (n *NatsConfig) GenerateConnString() string {
	var endpoints []string
	for _, v := range *n {
		endpoints = append(endpoints, fmt.Sprintf("%s:%s@%s", v.User, v.Password, v.Host))
	}
	return strings.Join(endpoints, ", ")
}

func (d *DatabaseConfig) GenerateConnString() string {
	var endpoints []string
	for _, v := range *d {
		endpoints = append(endpoints, fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", v.User, v.Password, v.Host, v.Port, v.Database))
	}
	return strings.Join(endpoints, ",")
}

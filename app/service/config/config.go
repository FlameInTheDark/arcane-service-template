package config

import (
	"fmt"
	"strings"
)

const (
	EtcdDatabase    = "/services/arcane/config/database"
	EtcdNats        = "/conf/nats/endpoints"
	EtcdEnvironment = "/services/arcane/config/environment"
)

type ConfigService struct {
	Environment EnvironmentConfig
	Nats        NatsConfig
	Database    DatabaseConfig
}

type EnvironmentConfig struct {
	Discord  string `json:"discord"`
	Database string `json:"database"`
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

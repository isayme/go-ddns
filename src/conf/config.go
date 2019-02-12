package conf

import (
	config "github.com/isayme/go-config"
	duration "github.com/isayme/go-duration"
)

// Logger ...
type Logger struct {
	Level string `json:"level"`
}

// DNSPod ...
type DNSPod struct {
	Domain    string            `json:"domain"`
	SubDomain string            `json:"subdomain"`
	Token     string            `json:"token"` // Token = Token ID + Token Value = ID,Value https://www.dnspod.cn/docs/info.html#common-parameters
	Interval  duration.Duration `json:"interval"`
}

// Config ...
type Config struct {
	Logger Logger `json:"logger"`

	DNSPod DNSPod `json:"dnspod"`

	IPURLs []string `json:"ipURLs"`
}

var globalConfig Config

// Get get config
func Get() *Config {
	config.Parse(&globalConfig)

	return &globalConfig
}

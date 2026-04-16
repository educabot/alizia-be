package config

import (
	"math"
	"time"

	bcfg "github.com/educabot/team-ai-toolkit/config"
)

type Config struct {
	bcfg.BaseConfig
	AzureOpenAIKey      string
	AzureOpenAIEndpoint string

	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration
}

func Load() *Config {
	base := bcfg.LoadBase()
	return &Config{
		BaseConfig:          base,
		AzureOpenAIKey:      bcfg.MustEnv("AZURE_OPENAI_API_KEY"),
		AzureOpenAIEndpoint: bcfg.MustEnv("AZURE_OPENAI_ENDPOINT"),

		DBMaxOpenConns:    boundedUintToInt(bcfg.GetEnvUint("DB_MAX_OPEN_CONNS", "25")),
		DBMaxIdleConns:    boundedUintToInt(bcfg.GetEnvUint("DB_MAX_IDLE_CONNS", "10")),
		DBConnMaxLifetime: time.Duration(boundedUintToInt(bcfg.GetEnvUint("DB_CONN_MAX_LIFETIME_MIN", "30"))) * time.Minute,
		DBConnMaxIdleTime: time.Duration(boundedUintToInt(bcfg.GetEnvUint("DB_CONN_MAX_IDLE_TIME_MIN", "5"))) * time.Minute,
	}
}

// boundedUintToInt clamps a uint value to math.MaxInt32 so the conversion to
// int is always safe regardless of platform word size.
func boundedUintToInt(v uint) int {
	if v > math.MaxInt32 {
		return math.MaxInt32
	}
	return int(v)
}

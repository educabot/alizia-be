package config

import (
	bcfg "github.com/educabot/team-ai-toolkit/config"
)

type Config struct {
	bcfg.BaseConfig
	AzureOpenAIKey      string
	AzureOpenAIEndpoint string
}

func Load() *Config {
	base := bcfg.LoadBase()
	return &Config{
		BaseConfig:          base,
		AzureOpenAIKey:      bcfg.MustEnv("AZURE_OPENAI_API_KEY"),
		AzureOpenAIEndpoint: bcfg.MustEnv("AZURE_OPENAI_ENDPOINT"),
	}
}

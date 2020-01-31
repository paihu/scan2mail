package scan2mail

import (
	toml "github.com/sioncojp/tomlssm"
)

// Config は設定を司る
type Config struct {
	BotToken string   `toml:"bot_token"`
	Domain   []string `toml:"valid_domain"`
	LogDir   string   `toml:"log_dir"`
}

// LoadToml はファイルからtomlを読む
func LoadToml(c, region string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(c, &config, region); err != nil {
		return nil, err
	}
	return &config, nil
}

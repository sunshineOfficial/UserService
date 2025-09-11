package config

import "github.com/sunshineOfficial/golib/config"

func Parse() (Settings, error) {
	var settings Settings
	return settings, config.Parse(&settings)
}

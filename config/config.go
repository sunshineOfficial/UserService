package config

type Settings struct {
	Port     int      `json:"port"`
	Database Database `json:"database"`
}

type Database struct {
	Postgres string `json:"postgres"`
}

func NewSettings() (Settings, error) {
	var settings Settings
	return settings, Parse(&settings)
}

package config

type Settings struct {
	Port     int      `json:"port"`
	Database Database `json:"database"`
	Kafka    Kafka    `json:"kafka"`
}

type Database struct {
	Postgres string `json:"postgres"`
}

type Kafka struct {
	Brokers []string `json:"brokers"`
	Topics  Topics   `json:"topics"`
}

type Topics struct {
	UserTickets string `json:"user_tickets"`
}

func NewSettings() (Settings, error) {
	var settings Settings
	return settings, Parse(&settings)
}

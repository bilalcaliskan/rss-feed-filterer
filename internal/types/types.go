package types

type Config struct {
	Repositories []Repository `yaml:"repositories"`
}

type Repository struct {
	Name                 string `yaml:"name"`
	Description          string `yaml:"description,omitempty"`
	RSSURL               string `yaml:"rssURL"`
	CheckIntervalMinutes int    `yaml:"checkIntervalMinutes"`
}

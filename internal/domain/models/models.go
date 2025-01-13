package models

type ServerConfig struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Debug   bool           `yaml:"debug"`
	Servers []ServerConfig `yaml:"servers"`
}

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

// TODO: WaybarOutput должен иметь метод для сериализации, который бы сам подставлял в json
// строку на основе внутренних данных + иконки - v0.1
const (
	EmptyInbox    = "<span rise='2000'>󰶈</span>"
	NonEmptyInbox = "<span color='#FF0000' rise='2000'>󰶍</span>"
)

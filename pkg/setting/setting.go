package setting

import "github.com/BurntSushi/toml"

type App struct {
	UserAgent       string
	RuntimeRootPath string
	LogSavePath     string
	LogSaveName     string
	LogFileExt      string
	TimeFormat      string
}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

type Settings struct {
	App   App
	Redis Redis
}

var GlobalSettings = &Settings{}

func Setup() error {
	// load configs from toml file
	_, err := toml.DecodeFile("conf/app.toml", GlobalSettings)
	if err != nil {
		return err
	}
	return nil
}

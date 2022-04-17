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

type Settings struct {
	App App
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

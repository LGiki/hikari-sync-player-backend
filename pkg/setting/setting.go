package setting

import "github.com/BurntSushi/toml"

type App struct {
	UserAgent       string
	Host            string
	Port            int
	RunningMode     string
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

var GlobalSettings = &Settings{
	App: App{
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36",
		Host:            "0.0.0.0",
		Port:            12321,
		RunningMode:     "release",
		RuntimeRootPath: "runtime/",
		LogSavePath:     "logs/",
		LogSaveName:     "log",
		LogFileExt:      "log",
		TimeFormat:      "20060102",
	},
	Redis: Redis{
		Host:        "127.0.0.1:6379",
		Password:    "",
		MaxIdle:     30,
		MaxActive:   30,
		IdleTimeout: 200,
	},
}

func Setup() error {
	// load configs from toml file
	_, err := toml.DecodeFile("conf/app.toml", GlobalSettings)
	if err != nil {
		return err
	}
	return nil
}

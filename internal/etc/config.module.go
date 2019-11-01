package etc

type Configuration struct {
	Name string
	Version string
	Db db
	Mqtt mqtt
	Web web
}

type db struct {
	Name string
	Url string
}

type mqtt struct {
	Broker string
	ClientId string `toml:"client-id"`
	SwitchId string `toml:"switch-id"`
	SwitchAppId string `toml:"switch-app-id"`
	SoilAppId string `toml:"soil-app-id"`
	SmokeAppId string `toml:"smoke-app-id"`
	Subscribe string
	Upstream string
}

type web struct {
	Listen string
}
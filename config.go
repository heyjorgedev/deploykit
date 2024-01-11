package deploykit

type DeploymentConfig struct {
	App struct {
		Name string `toml:"name"`
	} `toml:"app"`
}

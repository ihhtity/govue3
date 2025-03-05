package config

type Settings struct {
	Server struct {
		Port string `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"`

	Database struct {
		MySQL         string `mapstructure:"mysql"`
		Redis         string `mapstructure:"redis"`
		RedisPassword string `mapstructure:"redis_password"`
		MongoDB       string `mapstructure:"mongodb"`
	} `mapstructure:"database"`

	JWT struct {
		Secret string `mapstructure:"secret"`
		Expire int    `mapstructure:"expire"`
	} `mapstructure:"jwt"`

	Swagger struct {
		Enabled bool `mapstructure:"enabled"`
	} `mapstructure:"swagger"`
}

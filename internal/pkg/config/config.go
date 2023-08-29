package config

type Config struct {
	Name string `mapstructure:"Name"` //应用名称
	Host string `mapstructure:"Host"` //HTTP服务主机
	Port int    `mapstructure:"Port"` //HTTP端口

	MysqlConf struct {
		Host     string `mapstructure:"host"`     //
		Port     int    `mapstructure:"port"`     //
		Database string `mapstructure:"database"` //
		Username string `mapstructure:"username"` //
		Password string `mapstructure:"password"` //
		Charset  string `mapstructure:"charset"`  //https://mathiasbynens.be/notes/mysql-utf8mb4
	}

	RedisConf struct {
		Addr     string `mapstructure:"addr"` //
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"` //
		DB       int    `mapstructure:"db"`       //Redis库 建议默认0 因为Redis集群模式下没有16个库只有0号库
	}
}

var Conf = new(Config)

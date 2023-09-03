package config

type RedisConf struct {
	Addr     string `mapstructure:"addr"`     // redis实例主机地址
	Username string `mapstructure:"username"` // 用户名
	Password string `mapstructure:"password"` // 密码
	DB       int    `mapstructure:"db"`       // Redis库
	Network  string `mapstructure:"network"`  // 网络类型，tcp或unix，默认tcp
	ReadOnly bool   `mapstructure:"readOnly"` // 是否为只读实例
}

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

	//Redis配置
	RedisConf map[string]*RedisConf
}

var Conf = new(Config)

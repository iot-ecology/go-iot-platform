package config

type ServerConfig struct {
	NodeInfo     NodeInfo     `yaml:"node_info" json:"node_info"`
	RedisConfig  RedisConfig  `yaml:"redis_config" json:"redis_config"`
	MQConfig     MQConfig     `yaml:"mq_config" json:"mq_config"`
	InfluxConfig InfluxConfig `yaml:"influx_config" json:"influx_config"`
	MySQLConfig  MySQLConfig  `yaml:"mysql_config" json:"mysql_config"`
	MongoConfig  MongoConfig  `yaml:"mongo_config" json:"mongo_config"`
}
type MongoConfig struct {
	Host                   string `json:"host,omitempty" yaml:"host,omitempty"`
	Port                   int    `json:"port,omitempty" yaml:"port,omitempty"`
	Username               string `json:"username,omitempty" yaml:"username,omitempty"`
	Password               string `json:"password,omitempty" yaml:"password,omitempty"`
	Db                     string `json:"db,omitempty" yaml:"db,omitempty"`
	Collection             string `json:"collection,omitempty" yaml:"collection,omitempty"`
	WaringCollection       string `json:"waring_collection,omitempty" yaml:"waring_collection,omitempty"`
	ScriptWaringCollection string `json:"script_waring_collection,omitempty" yaml:"script_waring_collection,omitempty"`
}

type NodeInfo struct {
	Port int `json:"port,omitempty" yaml:"port,omitempty"`
}

type RedisConfig struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	Db       int    `json:"db,omitempty" yaml:"db,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type MQConfig struct {
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port,omitempty" yaml:"port,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

type InfluxConfig struct {
	Host   string `json:"host,omitempty" yaml:"host,omitempty"`
	Port   int    `json:"port,omitempty" yaml:"port,omitempty"`
	Token  string `json:"token,omitempty" yaml:"token,omitempty"`
	Org    string `json:"org,omitempty" yaml:"org,omitempty"`
	Bucket string `json:"bucket,omitempty" yaml:"bucket,omitempty"`
}
type MySQLConfig struct {
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Host     string `json:"host,omitempty" yaml:"host,omitempty"`
	Port     int    `json:"port" yaml:"port"`
	DBName   string `json:"dbname" yaml:"dbname"`
}

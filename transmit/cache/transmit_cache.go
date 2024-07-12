package cache

type MySQLTransmitCache struct {
	ID       string `json:"ID"`
	Host     string `json:"host"`
	Port     int    `json:"port" `
	Username string `json:"username" `
	Password string `json:"password" `
	Database string `json:"database" `
	Table    string `json:"table"`
	Script   string `json:"script"`
}

type MongoTransmitCache struct {
	ID         string `json:"ID"`
	Host       string `json:"host"`
	Port       int    `json:"port" `
	Username   string `json:"username" `
	Password   string `json:"password" `
	Database   string `json:"database" `
	Collection string `json:"collection"`
	Script     string `json:"script"`
}

type InfluxTransmitCache struct {
	ID          string `json:"ID"`
	Host        string `json:"host"`
	Port        int    `json:"port" `
	Token       string `json:"token"`
	Bucket      string `json:"bucket"`
	Org         string `json:"org"`
	Measurement string `json:"measurement"`
	Script      string `json:"script"`
}

type ClickhouseTransmitCache struct {
	ID       string `json:"ID"`
	Host     string `json:"host" `
	Port     int    `json:"port" `
	Username string `json:"username" `
	Password string `json:"password" `
	Database string `json:"database" `
	Table    string `json:"table"`
	Script   string `json:"script"`
}

type CassandraTransmitCache struct {
	ID       string `json:"ID"`
	Host     string `json:"host" `
	Port     int    `json:"port" `
	Username string `json:"username" `
	Password string `json:"password" `
	Database string `json:"database" `
	Table    string `json:"table"`
	Script   string `json:"script"`
}

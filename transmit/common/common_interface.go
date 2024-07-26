package common

type DataRowList struct {
	Time      int64     `json:"time"`       // 秒级时间戳
	DeviceUid string    `json:"device_uid"` // 是MqttClient的ID
	IdentificationCode string `json:"identification_code"` // 设备标识码
	DataRows  []DataRow `json:"data"`
	Nc        string    `json:"nc"`
}
type DataRow struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

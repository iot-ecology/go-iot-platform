package common

type DataRowList struct {
	Time      int64     `json:"time"`       // 秒级时间戳
	DeviceUid string    `json:"device_uid"` // 能够产生网络通讯的唯一编码
	IdentificationCode string `json:"identification_code"` // 设备标识码
	DataRows  []DataRow `json:"data"`
	Nc        string    `json:"nc"`
	Protocol string `json:"protocol"`
}
type DataRow struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

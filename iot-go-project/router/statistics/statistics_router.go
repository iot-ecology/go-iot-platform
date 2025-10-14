/*
Copyright 2024 - 2025 Zen HuiFer

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package statistics

import (
	"github.com/gin-gonic/gin"
	"igp/glob"
	"igp/servlet"
	"time"
)

type StatisticsApi struct {

}
type DeviceInfo struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Protocol  string    `gorm:"column:protocol"`
}

type DeviceStats struct {
	CreationDate string `json:"creation_date"`
	TotalDevices int    `json:"total_devices"`
	Protocol     string `json:"protocol"`
	ProtocolCount int   `json:"protocol_count"`
}

func (api *StatisticsApi) GetDeviceStats(c *gin.Context) {
	var stats []DeviceStats
	if err := glob.GDb.Table("device_infos").
		Select("DATE(created_at) AS creation_date, COUNT(*) AS total_devices, protocol, COUNT(*) AS protocol_count").
		Where("created_at IS NOT NULL").
		Group("creation_date, protocol").
		Order("creation_date DESC").
		Scan(&stats).Error; err != nil {
		servlet.Error(c, err.Error())
		return
	}
	servlet.Resp(c, stats)

}

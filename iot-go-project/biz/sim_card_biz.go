package biz

import (
	"context"
	"igp/glob"
	"igp/models"
	"igp/servlet"
	"strconv"
)

type SimCardBiz struct{}

func (biz *SimCardBiz) PageData(accessNumber string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var dashboards []models.SimCard

	db := glob.GDb
	if accessNumber != "" {
		db = db.Where("access_number like ?", "%"+accessNumber+"%")
	}
	db.Model(&models.SimCard{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&dashboards)

	pagination.Data = dashboards
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

var deviceInfoBiz = DeviceInfoBiz{}

func (biz *SimCardBiz) beforeCreate(card models.SimCard) {
	id := strconv.Itoa(int(card.ID))
	key := glob.SimCardExpireTime.String() + ":" + id
	glob.GRedis.SetNX(context.Background(), key, id, 0)
	glob.GRedis.ExpireAt(context.Background(), key, card.Expiration.AddDate(0, 0, -3))

}

func (biz *SimCardBiz) PageHistory(simCardId string, page, size int) (*servlet.PaginationQ, error) {
	var pagination servlet.PaginationQ
	var useHistories []models.SimUseHistory

	db := glob.GDb

	if simCardId != "" {
		db = db.Where("sim_id = ?", simCardId)
	}

	db.Model(&models.SimUseHistory{}).Count(&pagination.Total) // 计算总记录数
	offset := (page - 1) * size
	db.Offset(offset).Limit(size).Find(&useHistories)

	var resp []servlet.SimUseHistoryResp

	for _, dashboard := range useHistories {
		var id = deviceInfoBiz.FindById(dashboard.DeviceInfoId)

		historyResp := servlet.SimUseHistoryResp{
			ID:           dashboard.ID,
			SimId:        dashboard.SimId,
			DeviceInfoId: dashboard.DeviceInfoId,
			Description:  dashboard.Description,
		}
		if id != nil {
			historyResp.SN = id.SN
		}

		resp = append(resp, historyResp)
	}

	pagination.Data = resp
	pagination.Page = page
	pagination.Size = size

	return &pagination, nil
}

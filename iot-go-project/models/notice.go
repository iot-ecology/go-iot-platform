package models

import "gorm.io/gorm"

// DingDing 钉钉通知渠道  加签模式
type DingDing struct {
	gorm.Model
	Name        string `json:"name" structs:"name"` // 名称
	AccessToken string `json:"access_token" structs:"access_token"` 	// 访问令牌
	Secret      string `json:"secret" structs:"secret"` 	// 密钥
	Content     string `json:"content" structs:"content"` // 模板内容
}

// DingDingBindProduct 钉钉通知渠道绑定产品
type DingDingBindProduct struct {
	gorm.Model
	DingDingId int `json:"ding_ding_id" structs:"ding_ding_id"`
	ProductId  int `json:"product_id" structs:"product_id"`
}


// FeiShu 飞书通知渠道  加签模式
type FeiShu struct {
	gorm.Model
	Name        string `json:"name" structs:"name"`  // 名称
	AccessToken string `json:"access_token" structs:"access_token"` // 访问令牌
	Secret      string `json:"secret" structs:"secret"` 	// 密钥
	Content     string `json:"content" structs:"content"` // 模板内容
}

// FeiShuBindProduct 飞书通知渠道绑定产品
type FeiShuBindProduct struct {
	gorm.Model
	FeiShuId  int `json:"feishu_id" structs:"feishu_id"`
	ProductId int `json:"product_id" structs:"product_id"`
}

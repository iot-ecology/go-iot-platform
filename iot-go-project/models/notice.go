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

package models

import "gorm.io/gorm"

// DingDing 钉钉通知渠道  加签模式
type DingDing struct {
	gorm.Model `structs:"-"`
	Name        string `json:"name" structs:"name"` // 名称
	AccessToken string `json:"access_token" structs:"access_token"` 	// 访问令牌
	Secret      string `json:"secret" structs:"secret"` 	// 密钥
	Content     string `json:"content" structs:"content"` // 模板内容
}

// DingDingBindProduct 钉钉通知渠道绑定产品
type DingDingBindProduct struct {
	gorm.Model `structs:"-"`
	DingDingId int `json:"ding_ding_id" structs:"ding_ding_id"`
	ProductId  int `json:"product_id" structs:"product_id"`
}


// FeiShu 飞书通知渠道  加签模式
type FeiShu struct {
	gorm.Model `structs:"-"`
	Name        string `json:"name" structs:"name"`  // 名称
	AccessToken string `json:"access_token" structs:"access_token"` // 访问令牌
	Secret      string `json:"secret" structs:"secret"` 	// 密钥
	Content     string `json:"content" structs:"content"` // 模板内容
}

// FeiShuBindProduct 飞书通知渠道绑定产品
type FeiShuBindProduct struct {
	gorm.Model `structs:"-"`
	FeiShuId  int `json:"feishu_id" structs:"feishu_id"`
	ProductId int `json:"product_id" structs:"product_id"`
}

package models

import (
	"errors"
	"time"
)

type Validate interface {
	Validate() error
}

// Validate 验证SimCard实例是否满足业务规则
func (card *SimCard) Validate() error {
	// 验证接入号
	if card.AccessNumber == "" {
		return errors.New("接入号不能为空")
	}

	// 验证集成电路卡识别码
	if card.ICCID == "" {
		return errors.New("集成电路卡识别码不能为空")
	}

	// 验证国际移动用户识别码
	if card.IMSI == "" {
		return errors.New("国际移动用户识别码不能为空")
	}

	// 验证运营商名称
	if card.Operator == "" {
		return errors.New("运营商名称不能为空")
	}

	// 验证到期时间是否为未来时间
	if card.Expiration.Before(time.Now()) {
		return errors.New("到期时间必须是未来的日期")
	}

	// 如果所有验证都通过，则返回nil表示没有错误
	return nil
}

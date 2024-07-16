package test

import (
	"github.com/blinkbean/dingtalk"
	"iot-notice/models"
	"testing"
)

func TestMessageTemplate(t *testing.T) {
	temp := models.MessageTemplate{
		Content: "设备: {{device_name}} 信号: {{signal_name}} {{in_or_out}}阈值! \n当前值为{{signal_value}}{{unit}}\n阈值范围{{min" +
			"}}~{{max}}{{unit}}",
		DeviceName:   "测试设备",
		SignalId:     1,
		MqttClientId: 1,
		SignalName:   "温度",
		SignalValue:  10,
		Min:          0,
		Max:          3,
		Unit:         "℃",
		InOrOut:      1,
	}
	cli := dingtalk.InitDingTalkWithSecret("3b0263e8fcefee781f5fdcb19ecbc6178af530b530d259512de3500c69c94f8c", "SEC3d1a88943ae70a8dbfa1400fcba3ece5a0e79b33c15bf99803aa0f0f0708259b")
	cli.SendTextMessage(temp.Format())
}

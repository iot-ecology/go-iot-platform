package test

import (
	"github.com/blinkbean/dingtalk"
	"iot-notice/models"
	"log"
	"testing"

	"github.com/CatchZeng/feishu/pkg/feishu"
)

func TestMessageTemplate(t *testing.T) {
	temp := models.MessageTemplate{
		Content: "设备: {{device_name}} \n信号: {{signal_name}} {{in_or_out}}阈值! \n当前值: {{signal_value}}{{unit}}\n" +
			"阈值范围: {{min" +
			"}}~{{max}}{{unit}}",
		DeviceName:   "测试设备",
		SignalId:     1,
		MqttClientId: "1",
		SignalName:   "温度",
		SignalValue:  10,
		Min:          0,
		Max:          3,
		Unit:         "℃",
		InOrOut:      1,
	}
	cli := dingtalk.InitDingTalkWithSecret("3b0263e8fcefee781f5fdcb19ecbc6178af530b530d259512de3500c69c94f8c", "SEC3d1a88943ae70a8dbfa1400fcba3ece5a0e79b33c15bf99803aa0f0f0708259b")
	cli.SendTextMessage(temp.Format())


	token := "c081a7b8-06ab-46f0-b6c6-017652f5fbf9"
	secret := "IgZAYlSNO4HgIKabw4WrTd"

	client := feishu.NewClient(token, secret)
	text := feishu.NewText(temp.Format())
	line := []feishu.PostItem{text, }
	msg := feishu.NewPostMessage()

	msg.SetZHTitle("告警通知").
		AppendZHContent(line)
	req, resp, err := client.Send(msg)
	if err != nil {
		log.Print(err)
	}
	log.Print(resp)
	log.Print(req)
}


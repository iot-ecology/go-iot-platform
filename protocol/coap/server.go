package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-coap"
	"go.uber.org/zap"
	"log"
	"net"
)

func getUid(remoteAdd string) string {
	val := globalRedisClient.HGet(context.Background(), "coap_uid_f:"+globalConfig.NodeInfo.Name, remoteAdd).Val()
	return val
}

func storageUid(uid, remoteAdd string) {
	globalRedisClient.HSet(context.Background(), "coap_uid:"+globalConfig.NodeInfo.Name, uid, remoteAdd)
	globalRedisClient.HSet(context.Background(), "coap_uid_f:"+globalConfig.NodeInfo.Name, remoteAdd, uid)
}
func RemoveUid(remoteAdd string) {
	val := globalRedisClient.HGet(context.Background(), "coap_uid_f:"+globalConfig.NodeInfo.Name, remoteAdd).Val()
	if val == "" {
		return
	} else {
		globalRedisClient.HDel(context.Background(), "coap_uid:"+globalConfig.NodeInfo.Name, val)
		globalRedisClient.HDel(context.Background(), "coap_uid_f:"+globalConfig.NodeInfo.Name, remoteAdd)

	}
}

func Create(port int) {
	mux := coap.NewServeMux()
	mux.Handle("/auth", coap.FuncHandler(auth))                          //创建 "/auth"处理接口
	mux.Handle("/data", coap.FuncHandler(data))                          //创建 "/data"处理接口
	log.Fatal(coap.ListenAndServe("udp", fmt.Sprintf(":%d", port), mux)) //启动Server 端口为5683 这里为什么要用":5683"?搞不明白

}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DeviceId string `json:"device_id"`
}

func auth(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {

	if m.IsConfirmable() {

		auth := &Auth{}
		err := json.Unmarshal(m.Payload, auth)
		if err != nil {
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: m.MessageID,
				Token:     m.Token,
				Payload:   []byte("安全认证信息异常"),
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)
			return res
		}

		up, s := FindDeviceMappingUP(auth.DeviceId)

		if up != auth.Username || s != auth.Password {
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: m.MessageID,
				Token:     m.Token,
				Payload:   []byte("安全认证信息异常"),
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)
			return res
		}



		mc := globalRedisClient.HRandField(context.Background(), "coap_uid_f:"+globalConfig.NodeInfo.Name, -1).Val()

		if int64(len(mc)) <= globalConfig.NodeInfo.Size {
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: m.MessageID,
				Token:     m.Token,
				Payload:   []byte("安全认证成功"),
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)
			storageUid(auth.DeviceId, a.String())
			CoapMap[a.String()] = l
			return res
		}else {
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: m.MessageID,
				Token:     m.Token,
				Payload:   []byte("当前设备数量已达上线"),
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)
			return res
		}
	}
	return nil
}

var CoapMap = make(map[string]*net.UDPConn)

func data(l *net.UDPConn, a *net.UDPAddr, m *coap.Message) *coap.Message {
	if m.IsConfirmable() {
		uid := getUid(a.String())

		if uid != "" {
			coapMsg := CoapMessage{
				Uid:     uid,
				Message: string(m.Payload),
			}
			jsonData, err := json.Marshal(coapMsg)
			if err != nil {
				zap.S().Errorf("Error marshalling coap message to JSON: %v", err)
				res := &coap.Message{
					Type:      coap.Acknowledgement,
					Code:      coap.Content,
					MessageID: m.MessageID,
					Token:     m.Token,
					Payload:   []byte("数据序列化异常"),
				}
				res.SetOption(coap.ContentFormat, coap.TextPlain)
				return res
			}
			PushToQueue("pre_coap_handler", jsonData)
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: m.MessageID,
				Token:     m.Token,
				Payload:   []byte("数据处理成功"),
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)
			SetLastOpTime(a.String())
			DataHandlerCount()
			return res
		} else {
			res := &coap.Message{
				Type:      coap.Acknowledgement,
				Code:      coap.Content,
				MessageID: m.MessageID,
				Token:     m.Token,
				Payload:   []byte("你没有认证，请认证"),
			}
			res.SetOption(coap.ContentFormat, coap.TextPlain)
			return res
		}
	}
	return nil
}

type CoapMessage struct {
	Uid     string `json:"uid"`
	Message string `json:"message"`
}

func FindDeviceMappingUP(deviceId string) (string, string) {

	i := globalRedisClient.Exists(context.Background(), "auth:coap").Val()
	if i >= 0 {

		// todo: 从redis中根据deviceId获取用户名和密码
		val := globalRedisClient.HGet(context.Background(), "auth:coap", deviceId).Val()
		var auth Auth
		err := json.Unmarshal([]byte(val), &auth)
		if err != nil {
			zap.S().Errorf("Error unmarshalling auth to JSON: %v", err)
			return "", ""
		}
		return auth.Username, auth.Password
	}
	return "", ""
}

package influxdb2

import "testing"

func TestGet(t *testing.T) {
	host := "127.0.0.1"
	token := "i6XHSnNXeUoU3GoFXMm4qqrrgt69JKvQLqm0FCtnYG-rjb-nkDcry0pdwv4fpcXsSwi-mTGmAUTygkJtR-6CWA=="
	port := 8086
	client := GetInfluxDb(host, token, port, 1)
	client2 := GetInfluxDb(host, token, port, 1)
	if client == client2 {
		t.Logf("success")
	}
}

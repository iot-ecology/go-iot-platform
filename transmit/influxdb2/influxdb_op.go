package influxdb2

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"sync"
)

var influxClientMap = make(map[uint]influxdb2.Client)
var mu sync.Mutex

func GetInfluxDb(host, token string, port int, id uint) influxdb2.Client {
	mu.Lock()
	defer mu.Unlock()

	client, exists := influxClientMap[id]
	if exists {
		fmt.Printf("Reusing existing InfluxDB client for id: %d.\n", id)
		return client
	}

	client = influxdb2.NewClient(fmt.Sprintf("http://%s:%d", host, port), token)

	influxClientMap[id] = client
	fmt.Printf("New InfluxDB client for id %d has been initiated.\n", id)

	return client
}

package clickhouse

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {

	house1, _ := GetClickHouse(1, []string{"127.0.0.1:9000"}, "default", "default", "")
	house2, _ := GetClickHouse(1, []string{"127.0.0.1:9000"}, "default", "default", "")

	if house1 == house2 {
		t.Logf("success ")
	}
	query, err := house1.Query(context.Background(), "SELECT * FROM \"default\".\"sample\" LIMIT 300 OFFSET 0;")
	if err != nil {
		t.Error(err)
	} else {
		for query.Next() {
			var name string
			err := query.Scan(&name)
			if err != nil {
				t.Error(err)
			}
			t.Logf(" name:%s", name)
		}
	}
}

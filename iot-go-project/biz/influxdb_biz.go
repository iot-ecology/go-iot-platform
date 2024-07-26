package biz

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"igp/glob"
)

type InfluxdbBiz struct {
}

func (biz *InfluxdbBiz) QueryMeasurement(_measurement string) []string {

	query := fmt.Sprintf(`from(bucket: "%s")
              |> range(start: -1h, stop: now())
		      |> filter(fn: (r) => r._measurement =~ /^mqtt_%s/)
              |> keep(columns: ["_measurement"])
              |> group()
              |> distinct(column: "_measurement")
              |> limit(n: 10000)
              |> sort()`, glob.GConfig.InfluxConfig.Bucket, _measurement)

	queryAPI := glob.GInfluxdb.QueryAPI(glob.GConfig.InfluxConfig.Org)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		zap.S().Errorf("query error: %v", err)
		return []string{}
	}

	var res []string

	for result.Next() {
		record := result.Record()
		value := record.Value()
		res = append(res, value.(string))
	}
	return res
}

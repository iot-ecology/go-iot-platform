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

package biz

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"igp/glob"
)

type InfluxdbBiz struct {
}

func (biz *InfluxdbBiz) QueryMeasurement(_measurement , protocol string) []string {

	query := fmt.Sprintf(`from(bucket: "%s")
              |> range(start: -1h, stop: now())
		      |> filter(fn: (r) => r._measurement =~ /^%s_%s/)
              |> keep(columns: ["_measurement"])
              |> group()
              |> distinct(column: "_measurement")
              |> limit(n: 10000)
              |> sort()`, glob.GConfig.InfluxConfig.Bucket,protocol , _measurement)

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
